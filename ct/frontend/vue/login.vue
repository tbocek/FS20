<template>
    <div>
        <h1>Login</h1>
        <input type="text" name="email" v-model="email" placeholder="Email" />
        <input type="password" name="password" v-model="password" placeholder="Password" />
        <button type="button" v-on:click="login()">Login</button>
    </div>
</template>
<script>
    export default {
        data(){
            return {
                email : "",
                password : ""
            }
        },
        methods : {
            login() {
                if (this.password.length > 0) {
                    fetch('/auth', {
                        method: 'POST',
                        body: '{"username":"'+this.email+'", "password":"'+this.password+'"}'
                    }).then(response => {
                        if(response.status === 200) {
                            localStorage.setItem("token", response.headers.get("token"));
                            this.$router.push({ name: "main" });
                        }})
                }
            }
        }
    }
</script>