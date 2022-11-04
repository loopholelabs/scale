package io.loopholelabs.scale;

import io.loopholelabs.polyglot.*;
import io.loopholelabs.scale.generated.Context;

import org.teavm.interop.Import;
import org.teavm.interop.Export;
import org.teavm.interop.Address;

public class ScaleRunner {
  private static byte[] buffer = new byte[0];


	@Import(name = "next", module = "env")
	public static native int nextModule(int ptr, int size);

  @Export(name = "malloc")
  public static int malloc(int size) {
    buffer = new byte[size];
    return Address.ofData(buffer).toInt();
  }

  @Export(name = "resize")
  public static int resize(int size) {
    buffer = new byte[size];
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

      Context returnCtx = ScaleFunction.run(ctx);

      byte[] newdata = returnCtx.encode();
      long v = Address.ofData(newdata).toLong();
      v = (v << 32) | newdata.length;
      return v;

    } catch (DecodeException de) {
      return -1;  // Signal decoding error
    }
  }

	// Compile without, and this class won't be included. Hm.
	static void main(String... args) {}
}
