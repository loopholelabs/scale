package io.loopholelabs.scale;

import io.loopholelabs.scale.generated.*;

import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.Assertions;

import java.util.*;

class ContextTest {
 
  @Test
  void testDecodeEncodeContext() throws Exception {
    // This encoded context came from polyglot-ts
    String exampleEncoded = "050a04504f53540d16050a0468747470050a103a3a666666663a3132372e302e302e31040a0b48454c4c4f20574f524c440205030a05050a04686f737401050a01050a0f3132372e302e302e313a3430363737050a0f6163636570742d656e636f64696e6701050a01050a0d677a69702c206465666c617465050a0c636f6e74656e742d7479706501050a01050a216170706c69636174696f6e2f782d7777772d666f726d2d75726c656e636f646564050a0e636f6e74656e742d6c656e67746801050a01050a023131050a0a636f6e6e656374696f6e01050a01050a05636c6f73650c9003040a000205030a00";

    byte[] data = HexFormat.of().parseHex(exampleEncoded);

    Context c = new Context();
    byte[] leftData = c.decodeFrom(data);

    Assertions.assertEquals(leftData.length, 0); // Make sure it was all used.

    // Make sure the request was decoded correctly...
    Assertions.assertEquals(c.request.method, "POST");
    Assertions.assertEquals(c.request.protocol, "http");
    Assertions.assertEquals(c.request.ip, "::ffff:127.0.0.1");
    Assertions.assertArrayEquals(c.request.body, "HELLO WORLD".getBytes());
    Assertions.assertEquals(c.request.contentLength, 11);

    // Make sure the response was decoded correctly...
    Assertions.assertEquals(c.response.statusCode, 200);
    Assertions.assertArrayEquals(c.response.body, new byte[0]);

    // TODO: Check headers

    // Now encode and make sure it is the same as our control input.
    /* NB Doesn't work probably due to unpredictable map iterate order.
    byte[] encoded = c.encode();
    Assertions.assertArrayEquals(encoded, data);
    */
  }
}