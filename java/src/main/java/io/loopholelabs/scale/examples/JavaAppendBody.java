package io.loopholelabs.scale.examples;

import io.loopholelabs.scale.generated.*;

public class JavaAppendBody {

  /**
   * Simple scale function example
   * 
   * @param c
   * @return
   */
  public static Context run(Context c) {
    c.response.body = (new String(c.response.body) + " JavaAppendBody").getBytes();
    return c.next();
  }

}