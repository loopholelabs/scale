/*
	Copyright 2022 Loophole Labs
	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at
		   http://www.apache.org/licenses/LICENSE-2.0
	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

package polyglot;

public class Buffer {
  private byte[] buffer;

  public Buffer() {
    this.buffer = new byte[0];
  }

  public void Reset() {
    this.buffer = new byte[0];
  }

  public int Write(byte data) {
    byte[] someData = new byte[1];
    someData[0] = data;
    return this.Write(someData);
  }

  public int Write(byte[] data) {
    if (data.length==0) {
      return 0;
    }
    byte[] newBuffer = new byte[this.buffer.length + data.length];
    System.arraycopy(this.buffer, 0, newBuffer, 0, this.buffer.length);
    System.arraycopy(data, 0, newBuffer, this.buffer.length, data.length);
    this.buffer = newBuffer;
    return data.length;
  }

  public byte[] Bytes() {
    return this.buffer;
  }

  public int Len() {
    return this.buffer.length;
  }
}