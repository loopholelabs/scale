package polyglot;

import polyglot.generated.*;

import java.util.*;

public class Demo {

  public static String exampleEncoded = "050a04504f53540d16050a0468747470050a103a3a666666663a3132372e302e302e31040a0b48454c4c4f20574f524c440205030a05050a04686f737401050a01050a0f3132372e302e302e313a3430363737050a0f6163636570742d656e636f64696e6701050a01050a0d677a69702c206465666c617465050a0c636f6e74656e742d7479706501050a01050a216170706c69636174696f6e2f782d7777772d666f726d2d75726c656e636f646564050a0e636f6e74656e742d6c656e67746801050a01050a023131050a0a636f6e6e656374696f6e01050a01050a05636c6f73650c9003040a000205030a00";
  
  public static void main(String[] args) throws Exception {
    System.out.println("Starting test example...\n" + exampleEncoded);

    byte[] data = HexFormat.of().parseHex(exampleEncoded);

    // Now decode it

    Context c = new Context();
    data = c.decodeFrom(data);

    System.out.println(c);

    // Encode it and check...

    byte[] d = c.encode();
    System.out.println("Encoded\n" + HexFormat.of().formatHex(d));
  }
}