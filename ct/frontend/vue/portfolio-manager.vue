<template>
    <div>
        <p v-for="item in userCoinPortfolio">
            <input v-model="item.symbol" placeholder="BTC">
            <input v-model="item.amount" placeholder="0">
            {{item.amount * find(item.symbol)}}
            <button @click="remove(item.symbol)">remove</button>
        </p>
        <p>Total: {{total | round}} USD</p>
        <select v-model="selected">
            <option disabled value="">Please select one</option>
            <option v-for="item in coinmarketcapData" :value="item.symbol">
                {{ item.name }}
            </option>
        </select>
        <button @click="userCoinPortfolio.push({'symbol' : selected, 'amount': '0'})">add</button>
        <button @click="logout()">logout</button>
    </div>
</template>
<script>
    export default {
        data: function () {
            return {
                coinmarketcapData: [],
                userCoinPortfolio: [],
                selected: '',
                token: ''
            }
        },
        mounted: function () {
            fetch('/proxy')
                .then(response => {
                    if(response.status !== 200) {
                        throw Error("could not fetch proxy: " + response);
                    }
                    return response.json()
                }).then(body => {
                    this.coinmarketcapData = body.data
                });
            this.token = localStorage.getItem("token");

            fetch('/coinservice/portfolio', {
                method: 'GET',
                headers: new Headers({
                    'Authorization': 'Bearer ' + this.token,
                    'Content-Type': 'application/application/json'
                }),
            }).then(response => {
                if(response.status === 401) {
                    this.$router.push({ name: "login" });
                } else if(response.status !== 200) {
                    throw Error("could not fetch get data: " + response);
                }
                return response.json()})
                .then(body => {
                    this.userCoinPortfolio = body;
                });

        },
        methods: {
            logout: function () {
                localStorage.setItem("token", "");
                this.token = "";
                this.$router.push({ name: "login" });
            },
            find: function (symbol) {
                if (this.coinmarketcapData.length > 0) {
                    let found = this.coinmarketcapData.find(function (element) {
                        return element.symbol === symbol
                    });
                    return parseFloat(found.quote.USD.price)
                } else return 0;
            },
            remove: function (symbol) {
                this.userCoinPortfolio = this.userCoinPortfolio.filter(item => item.symbol !== symbol)
            }
        },
        computed: {
            total: function () {
                return this.userCoinPortfolio.reduce((total, item) => {
                    return total + (item.amount * this.find(item.symbol));
                }, 0)
            }
        },
        watch: {
            userCoinPortfolio: {
                handler: function () {
                    fetch('/coinservice/portfolio', {
                        method: 'POST',
                        headers: new Headers({
                            'Authorization': 'Bearer ' + this.token,
                            'Content-Type': 'application/application/json'
                        }),
                        body: JSON.stringify(this.userCoinPortfolio)
                    }).then(response => {
                        if(response.status === 401) {
                            this.$router.push({name: "login"});
                        } else if(response.status !== 200) {
                            throw Error("could not fetch post data: " + response);
                        }
                    })
                }, deep: true
            }
        },
        filters: {
            round: function (value) {
                return Math.round(value)
            }
        }
    }
</script>
