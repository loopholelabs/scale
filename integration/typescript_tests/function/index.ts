import * as signature from "signature";

export function example(ctx?: signature.ModelWithAllFieldTypes): signature.ModelWithAllFieldTypes | undefined {
    console.log("This is a Typescript Function");
    if (typeof ctx !== "undefined") {
        ctx.stringField = "This is a Typescript Function"
    }
    return signature.Next(ctx);
}