package polyglot.generated;

import polyglot.*;

import java.util.*;

public class Request {
  public String method;
  public long contentLength;
  public String protocol;
  public String ip;
  public byte[] body;
  public HashMap<String, String[]> headers;

  public byte[] encode() {
    Buffer b = new Buffer();
    Encode.encodeString(b, method);
    Encode.encodeInt64(b, contentLength);
    Encode.encodeString(b, protocol);
    Encode.encodeString(b, ip);
    Encode.encodeBytes(b, body);
    Encode.encodeMap(b, headers.size(), Encode.StringKind, Encode.AnyKind);
    Iterator<String> i = headers.keySet().iterator();
    while(i.hasNext()) {
      String k = i.next();
      Encode.encodeString(b, k);
      String[] v = headers.get(k);
      Encode.encodeSlice(b, v.length, Encode.StringKind);
      for(int j=0;j<v.length;j++) {
        // Encode string array
        Encode.encodeString(b, v[j]);
      }
    }
    return b.Bytes();
  }

  public byte[] decodeFrom(byte[] data) throws DecodeException {
    this.method = Decode.decodeString(data);
    data = Decode.skipString(data);
    this.contentLength = Decode.decodeInt64(data);
    data = Decode.skipInt64(data);
    this.protocol = Decode.decodeString(data);
    data = Decode.skipString(data);
    this.ip = Decode.decodeString(data);
    data = Decode.skipString(data);
    this.body = Decode.decodeUint8Array(data);
    data = Decode.skipUint8Array(data);

    headers = new HashMap<String, String[]>();

    Decode.MapInfo mi = Decode.decodeMap(data);
    data = Decode.skipMap(data);

    for(int i=0;i<mi.size;i++) {
      // Decode each item...
      String k = Decode.decodeString(data);
      data = Decode.skipString(data);

      Decode.SliceInfo si = Decode.decodeSlice(data);
      data = Decode.skipSlice(data);

      String[] vals = new String[si.size];
      for (int j=0;j<si.size;j++) {
        vals[j] = Decode.decodeString(data);
        data = Decode.skipString(data);
      }
      headers.put(k, vals);
    }

    return data;
  }

  public String toString() {
    String s = "Request[Method=" + method + ", contentLength=" + contentLength + ", protocol=" + protocol + ", IP=" + ip + "]\n" +
           " Body=" + new String(body) + "\n";
    Iterator<String> i = headers.keySet().iterator();
    while(i.hasNext()) {
      String k = i.next();
      String[] v = headers.get(k);
      for(int j=0;j<v.length;j++) {
        s = s + " Header " + k + " = " + v[j] + "\n";
      }
    }
    return s;
  }
}