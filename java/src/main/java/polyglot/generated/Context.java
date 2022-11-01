package polyglot.generated;

import polyglot.*;

public class Context {

  public Request request;
  public Response response;

  public Context() {

  }
  
  public byte[] decodeFrom(byte[] data) throws DecodeException {
    this.request = new Request();
    data = this.request.decodeFrom(data);
    this.response = new Response();
    data = this.response.decodeFrom(data);
    return data;
  }

  // Encode to bytes...
  public byte[] encode() {
    byte[] reqData = request.encode();
    byte[] respData = response.encode();
    byte[] data = new byte[reqData.length + respData.length];
    for(int i=0;i<reqData.length;i++) {data[i] = reqData[i];}
    for(int i=0;i<respData.length;i++) {data[i+reqData.length] = respData[i];}
    return data;
  }

  public String toString() {
    return "Context\n" + request + "\n" + response;
  }
}