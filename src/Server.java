import java.net.*;

class Server
{
    public static void main(String args[]) throws Exception
    {
        DatagramSocket serverSocket = new DatagramSocket(8081);
        byte[] receiveData = new byte[1024];
        while(true)
        {
            DatagramPacket receivePacket = new DatagramPacket(receiveData, receiveData.length);
            serverSocket.receive(receivePacket);
            String s = new String( receivePacket.getData(), 0, receivePacket.getLength());
            System.out.println("Message Received: " + s);
        }
    }
}
