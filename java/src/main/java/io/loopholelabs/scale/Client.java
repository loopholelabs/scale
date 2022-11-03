package io.loopholelabs.scale;

import io.loopholelabs.polyglot.*;
import io.loopholelabs.scale.generated.Context;

import org.teavm.interop.Import;
import org.teavm.interop.Export;
import org.teavm.interop.Address;

public class Client {

	@Import(name = "next", module = "env")
	public static native int nextModule(int ptr, int size);

  @Export(name = "malloc")
  public static int malloc(int size) {
    byte[] buffer = new byte[size];
    return Address.ofData(buffer).toInt();
  }

  @Export(name = "run")
  public static long run(int ptr, int size) {

    Address a = Address.fromInt(ptr);

    byte[] data = new byte[size];

    for (int i=0;i<size;i++) {
      data[i] = a.getByte();
      a = a.add(1);
    }

    try {
      Context ctx = new Context();
      ctx.decodeFrom(data);

      System.out.println("Got context into JAVA as " + ctx);

      ctx.response.statusCode = 299;
      ctx.response.body = "Hello from java!".getBytes();
      String[] vals = new String[1];
      vals[0] = "Hello world";
      ctx.response.headers.put("JAVA", vals);

      // Now encode it, and return the ptr/len

      byte[] newdata = ctx.encode();
      long v = Address.ofData(newdata).toLong();
      v = (v << 32) | newdata.length;
      return v;

    } catch (DecodeException de) {
      return 999;
    }
  }

	// Compile without, and this class won't be included. Hm.
	static void main(String... args) {}
}
