// setup jest
import 'jest';
import {
    NeverPullPolicy,
    defaultOrganization,
    New,
    WithCacheDirectory,
    WithPullPolicy,
    WithAPIKey, WithBaseURL, WithOrganization, computeSHA256
} from './index';
import {ScaleFunc, V1Alpha, Go} from "@loopholelabs/scalefile";
import {Storage} from "../storage";

const testingApiBaseUrl = "https://api.dev.scale.sh/v1"

/**
 * @jest-environment node
 */
test('TestPulldownCache', async () => {
    const st = new Storage("testDirectory");

    const scaleFunc = new ScaleFunc(
        V1Alpha,
        "TestFunction",
        "1",
        "signature1",
        Go,
        Buffer.from("Hello world"),
    );

    const func = "TestFunction"
    const tag = "1"
    const hash = await computeSHA256(scaleFunc.Function);
    st.Put(func, tag, defaultOrganization, hash, scaleFunc)

    const newScaleFunc = await New(
        func,
        tag,
        WithCacheDirectory("testDirectory"),
        WithPullPolicy(NeverPullPolicy),
    );

    expect(newScaleFunc.Version).toBe(scaleFunc.Version);
    expect(newScaleFunc.Name).toBe(scaleFunc.Name);
    expect(newScaleFunc.Tag).toBe(scaleFunc.Tag);
    expect(newScaleFunc.Signature).toBe(scaleFunc.Signature);
    expect(newScaleFunc.Language).toBe(scaleFunc.Language);
    expect(newScaleFunc.Function).toEqual(scaleFunc.Function);
    expect(newScaleFunc.Size).toBe(61);
});

test('TestRegistryDownload', async () => {
   // get api key from environment
    const apiKey = process.env.SCALE_API_KEY;
    if (apiKey == undefined) {
        console.log("SCALE_API_KEY not set, skipping test")
        return
    }

    const scaleFunc = await New(
        "TestRegistryDownload",
        "1",
        WithAPIKey(apiKey),
        WithBaseURL(testingApiBaseUrl),
        WithOrganization("alex"),
    );

    expect(scaleFunc.Version).toBe("v1alpha");
    expect(scaleFunc.Name).toBe("TestRegistryDownload");
    expect(scaleFunc.Tag).toBe("1");
    expect(scaleFunc.Signature).toBe("signature1");
    expect(scaleFunc.Language).toBe("go");
});
