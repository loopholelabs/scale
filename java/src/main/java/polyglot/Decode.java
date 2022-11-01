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

import java.util.*;

public class Decode {

  public static DecodeException InvalidSlice = new DecodeException("Invalid slice encoding");
  public static DecodeException InvalidMap = new DecodeException("Invalid map encoding");
  public static DecodeException InvalidBytes = new DecodeException("Invalid bytes encoding");
  public static DecodeException InvalidString = new DecodeException("Invalid string encoding");
  public static DecodeException InvalidError = new DecodeException("Invalid error encoding");
  public static DecodeException InvalidBool = new DecodeException("Invalid bool encoding");
  public static DecodeException InvalidUint8 = new DecodeException("Invalid uint8 encoding");
  public static DecodeException InvalidUint16 = new DecodeException("Invalid uint16 encoding");
  public static DecodeException InvalidUint32 = new DecodeException("Invalid uint32 encoding");
  public static DecodeException InvalidUint64 = new DecodeException("Invalid uint64 encoding");
  public static DecodeException InvalidInt32 = new DecodeException("Invalid int32 encoding");
  public static DecodeException InvalidInt64 = new DecodeException("Invalid int64 encoding");
  public static DecodeException InvalidFloat32 = new DecodeException("Invalid float32 encoding");
  public static DecodeException InvalidFloat64 = new DecodeException("Invalid float64 encoding");


  public static class MapInfo {
    public int size;
    public byte keyKind;
    public byte valueKind;
  }

  public static MapInfo decodeMap(byte[] data) throws DecodeException {
    if (data[0] != Encode.MapKind) {
      throw InvalidMap;
    }
    
    MapInfo mi = new MapInfo();
    mi.keyKind = data[1];
    mi.valueKind = data[2];
    data = Arrays.copyOfRange(data, 3, data.length);    
    mi.size = decodeUint32(data);
    return mi;
  }

  public static byte[] skipMap(byte[] data) throws DecodeException {
    data = Arrays.copyOfRange(data, 3, data.length);    
    return skipUint32(data);
  }

  public static class SliceInfo {
    public int size;
    public byte valueKind;
  }

  public static SliceInfo decodeSlice(byte[] data) throws DecodeException {
    if (data[0] != Encode.SliceKind) {
      throw InvalidSlice;
    }
    
    SliceInfo si = new SliceInfo();
    si.valueKind = data[1];
    data = Arrays.copyOfRange(data, 2, data.length);    
    si.size = decodeUint32(data);
    return si;
  }

  public static byte[] skipSlice(byte[] data) throws DecodeException {
    data = Arrays.copyOfRange(data, 2, data.length);    
    return skipUint32(data);
  }

  public static String decodeString(byte[] data) throws DecodeException {
    if (data[0] != Encode.StringKind) {
      throw InvalidString;
    }
    data = Arrays.copyOfRange(data, 1, data.length);
    int len = decodeUint32(data);
    data = skipUint32(data);
    return new String(data, 0, len);
  }

  public static byte[] skipString(byte[] data) throws DecodeException {
    if (data[0] != Encode.StringKind) {
      throw InvalidString;
    }
    data = Arrays.copyOfRange(data, 1, data.length);
    int len = decodeUint32(data);
    data = skipUint32(data);
    data = Arrays.copyOfRange(data, len, data.length);
    return data;
  }

  public static int decodeUint32(byte[] data) throws DecodeException {
    if (data[0] != Encode.Uint32Kind) {
      throw InvalidUint32;
    }
    data = Arrays.copyOfRange(data, 1, data.length);

    int x = 0;
    int s = 0;

    for(int i=0; i < Encode.VarIntLen32; i++) {
      int cb = ((int)data[i]) & 0x7f;
      // Check if msb is set signifying a continuation byte

      if ((data[i] & Encode.continuation) == 0) {
        if (i >= Encode.VarIntLen32 && cb > 1) {
          throw InvalidUint32;
        }
        // End of varint, add the last bits and advance the buffer
        return x | (cb << s);
      }
      // Add the lower 7 bits to the result and continue to the next byte
      x |= cb << s;
      s += 7;
    }
    throw InvalidUint32;
  }

  public static byte[] skipUint32(byte[] data) throws DecodeException {
    if (data[0] != Encode.Uint32Kind) {
      throw InvalidUint32;
    }

    data = Arrays.copyOfRange(data, 1, data.length);

    for(int i=0; i < Encode.VarIntLen32; i++) {
      int cb = ((int)data[i]) & 0xff;
      // Check if msb is set signifying a continuation byte
      if ((data[i] & Encode.continuation) == 0) {
        if (i >= Encode.VarIntLen32 && cb > 1) {
          throw InvalidUint32;
        }
        // Get rid of the first i bytes from the buffer and return it...
        return Arrays.copyOfRange(data, i + 1, data.length);
      }
    }
    throw InvalidUint32;
  }

  public static int decodeInt32(byte[] data) throws DecodeException {
    if (data[0] != Encode.Int32Kind) {
      throw InvalidInt32;
    }

    data = Arrays.copyOfRange(data, 1, data.length);

    long ux = 0;
    long s = 0;

    for (int i = 0; i < Encode.VarIntLen32; i++) {
      int cb = ((int)data[i] & 0x7f);

      // Check if msb is set signifying a continuation byte
      if ((data[i] & Encode.continuation) == 0) {
        if (i >= Encode.VarIntLen32 && cb > 1) {
          throw InvalidInt32;
        }
        // End of varint, add the last bits and cast to signed integer
        long x = (ux | (cb <<s )) >> 1;

        if ((ux & 1) != 0) {
          x = ~x;
        }

        return (int)x;
      }
      ux |= cb << s;
      s += 7;
    }

    throw InvalidInt32;
  }

  public static byte[] skipInt32(byte[] data) throws DecodeException {
    if (data[0] != Encode.Int32Kind) {
      throw InvalidInt32;
    }

    data = Arrays.copyOfRange(data, 1, data.length);

    for (int i = 0; i < Encode.VarIntLen32; i++) {
      int cb = ((int)data[i] & 0x7f);

      // Check if msb is set signifying a continuation byte
      if ((data[i] & Encode.continuation) == 0) {
        if (i >= Encode.VarIntLen32 && cb > 1) {
          throw InvalidInt64;
        }
        return Arrays.copyOfRange(data, i + 1, data.length);
      }
    }

    throw InvalidInt32;
  }

  public static long decodeInt64(byte[] data) throws DecodeException {
    if (data[0] != Encode.Int64Kind) {
      throw InvalidInt64;
    }

    data = Arrays.copyOfRange(data, 1, data.length);

    long ux = 0;
    long s = 0;

    for (int i = 0; i < Encode.VarIntLen64; i++) {
      int cb = ((int)data[i] & 0x7f);

      // Check if msb is set signifying a continuation byte
      if ((data[i] & Encode.continuation) == 0) {
        if (i >= Encode.VarIntLen64 && cb > 1) {
          throw InvalidInt64;
        }
        // End of varint, add the last bits and cast to signed integer
        long x = (ux | (cb <<s )) >> 1;
        // TODO: Flip the bits if the sign bit is set, cope with overflow above.
/*
        if ((ux & 1) != 0) {
          x = ^x;
        }
*/
        return x;
      }
      ux |= cb << s;
      s += 7;
    }

    throw InvalidInt64;
  }

  public static byte[] skipInt64(byte[] data) throws DecodeException {
    if (data[0] != Encode.Int64Kind) {
      throw InvalidInt64;
    }

    data = Arrays.copyOfRange(data, 1, data.length);

    for (int i = 0; i < Encode.VarIntLen64; i++) {
      int cb = ((int)data[i] & 0x7f);

      // Check if msb is set signifying a continuation byte
      if ((data[i] & Encode.continuation) == 0) {
        if (i >= Encode.VarIntLen64 && cb > 1) {
          throw InvalidInt64;
        }
        return Arrays.copyOfRange(data, i + 1, data.length);
      }
    }

    throw InvalidInt64;
  }

  public static byte[] decodeUint8Array(byte[] data) throws DecodeException {
    if (data[0]!=Encode.BytesKind) {
      throw InvalidBytes;
    }
    data = Arrays.copyOfRange(data, 1, data.length);

    int len = decodeUint32(data);
    data = skipUint32(data);
    return Arrays.copyOfRange(data, 0, len);
  }

  public static byte[] skipUint8Array(byte[] data) throws DecodeException {
    if (data[0]!=Encode.BytesKind) {
      throw InvalidBytes;
    }
    data = Arrays.copyOfRange(data, 1, data.length);

    int len = decodeUint32(data);
    data = skipUint32(data);
    return Arrays.copyOfRange(data, len, data.length);
  }

}