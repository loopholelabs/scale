package io.loopholelabs.scale.generated;

import io.loopholelabs.polyglot.*;

import java.util.*;

public class Response {
  public int statusCode;
  public byte[] body;
  public HashMap<String, String[]> headers;

  public byte[] encode() {
    Buffer b = new Buffer();
    Encode.encodeInt32(b, statusCode);
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
    this.statusCode = Decode.decodeInt32(data);
    data = Decode.skipInt32(data);
    this.body = Decode.decodeBytes(data);
    data = Decode.skipBytes(data);

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
    String s = "Response[Status=" + statusCode + "\n" +
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