package io.loopholelabs.scale.examples;

import io.loopholelabs.scale.generated.*;

public class JavaEndpoint {

  /**
   * Simple scale function example
   * 
   * @param c
   * @return
   */
  public static Context run(Context c) {
    c.response.statusCode = 299;
    c.response.body = "Hello from java!".getBytes();
    String[] vals = new String[1];
    vals[0] = "Hello world";
    c.response.headers.put("JAVA", vals);
    return c;
  }

}