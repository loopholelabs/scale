package io.loopholelabs.scale.examples;

import io.loopholelabs.scale.generated.*;

public class JavaHeaders {

  /**
   * Simple scale function example
   * 
   * @param c
   * @return
   */
  public static Context run(Context c) {
    String[] vals = new String[1];
    vals[0] = "Hello";
    c.response.headers.put("JAVAheader", vals);
    return c.next();
  }

}