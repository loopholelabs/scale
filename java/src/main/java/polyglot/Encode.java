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

public class Encode {
  public static long continuation = 0x80;
  public static int VarIntLen16 = 3;
  public static int VarIntLen32 = 5;
  public static int VarIntLen64 = 10;

  public static byte NilKind = 0;
  public static byte SliceKind = 1;
  public static byte MapKind = 2;
  public static byte AnyKind = 3;
  public static byte BytesKind = 4;
  public static byte StringKind = 5;
  public static byte ErrorKind = 6;
  public static byte BoolKind = 7;
  public static byte Uint8Kind = 8;
  public static byte Uint16Kind = 9;
  public static byte Uint32Kind = 10;
  public static byte Uint64Kind = 11;
  public static byte Int32Kind = 12;
  public static byte Int64Kind = 13;
  public static byte Float32Kind = 14;
  public static byte Float64Kind = 15;

  public static byte FalseBool = 0;
  public static byte TrueBool = 1;

  public static void encodeNil(Buffer b) {
    b.Write(NilKind);
  }

  public static void encodeMap(Buffer b, int size, byte keyKind, byte valueKind) {
    b.Write(MapKind);
    b.Write(keyKind);
    b.Write(valueKind);
    encodeUint32(b, size);
  }

  public static void encodeSlice(Buffer b, int size, byte kind) {
    b.Write(SliceKind);
    b.Write(kind);
    encodeUint32(b, size);
  }

  public static void encodeBytes(Buffer b, byte[] value) {
    b.Write(BytesKind);
    encodeUint32(b, value.length);
    b.Write(value);
  }

  public static void encodeString(Buffer b, String value) {
    b.Write(StringKind);
    encodeUint32(b, value.getBytes().length);
    b.Write(value.getBytes());
  }

  public static void encodeError(Buffer b, String value) {
    b.Write(ErrorKind);
    encodeString(b, value);   
  }

  public static void encodeBool(Buffer b, boolean value) {
    b.Write(BoolKind);
    if (value) {
      b.Write(TrueBool);
    } else {
      b.Write(FalseBool);
    }
  }

  public static void encodeUint8(Buffer b, byte value) {
    b.Write(Uint8Kind);
    b.Write(value);
  }

  public static void encodeUint16(Buffer b, short value) {
    b.Write(Uint16Kind);
    while(value >= continuation) {
      b.Write((byte)(value | continuation));
      value >>>= 7;
    }
    b.Write((byte)value);
  }

  public static void encodeUint32(Buffer b, int value) {
    b.Write(Uint32Kind);
    while(value >= continuation) {
      b.Write((byte)(value | continuation));
      value >>>= 7;
    }
    b.Write((byte)value);
  }

  public static void encodeUint64(Buffer b, long value) {
    b.Write(Uint64Kind);
    while(value >= continuation) {
      b.Write((byte)(value | continuation));
      value >>>= 7;
    }
    b.Write((byte)value);
  }

  public static void encodeInt32(Buffer b, int value) {
    b.Write(Int32Kind);

    // Shift the value to the left by 1 bit, then flip the bits if the value is negative.
    long castValue = value << 1;  
    if (value < 0) {
      castValue = ~castValue;
    }

    while(castValue >= continuation) {
      // Append the lower 7 bits of the value, then shift the value to the right by 7 bits.
      b.Write((byte)(castValue | continuation));
      castValue >>= 7;
    }
    b.Write((byte)castValue);
  }

  public static void encodeInt64(Buffer b, long value) {
    b.Write(Int64Kind);

    // Shift the value to the left by 1 bit, then flip the bits if the value is negative.
    long castValue = value << 1;
    if (value < 0) {
      castValue = ~castValue;
    }

    while(castValue >= continuation) {
      // Append the lower 7 bits of the value, then shift the value to the right by 7 bits.
      b.Write((byte)(castValue | continuation));
      castValue >>= 7;
    }
    b.Write((byte)castValue);
  }

  public static void encodeFloat32(Buffer b, float value) {
    // TODO
  }

  public static void encodeFloat64(Buffer b, double value) {
    // TODO
  }

}