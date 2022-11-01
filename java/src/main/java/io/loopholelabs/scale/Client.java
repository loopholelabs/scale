package io.loopholelabs.scale;

import org.teavm.interop.Import;
import org.teavm.interop.Export;

public class Client {

	public static int counter = 0;

	@Import(name = "inc_by", module = "teavm_unchained")
	public static native int incBy();

	@Export(name = "inc")
	public static int inc() {
		counter += incBy();
		System.out.println("Current counter (VM): " + counter);
		return counter;
	}

  @Export(name = "run")
  public static int run(int ptr, int size) {
    System.out.println("Run was called with " + ptr + ", " + size);
    return ptr * 2;
  }

	// Compile without, and this class won't be included. Hm.
	static void main(String... args) {}
}
