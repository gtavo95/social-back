


# Create a Social Post

This is a GO backend project. You can consume it using the next TS

First create some Types


```ts
export type Params = {
	tone: string;
	words: string;
	hashtags: boolean;
	emojis: boolean
	network: string
	context: boolean
	posts: string
	url: string
}

export type Meeting = {
	startTime: string
	endTime: string
	link: string
}

```

## From the server side create and send the post.

```ts

let params: Params = {
    tone: "polite",
    words: "20",
    hashtags: true,
    emojis: true,
    network: "x",
    url: "*ecomerce page*"
}

let meeting: Meeting = {
    startTime: "MonthName/Day HH:MM.",
    endTime: "MonthName/Day HH:MM.",
    link: "*meeting link*"
}

let prompt = "*The user prompt goes here*"

const data = JSON.stringify({
    prompt: prompt,
    params: params,
    meeting: meeting
})

const res = await fetch("https://social-back-531344799107.us-central1.run.app/text", {
    method: "POST",
    body: data,
})

if (res.ok) {
    const response = await res.json()
    console.log(response)
    return response
}


```
Now read the incoming data on the client side

```ts
if (
    result.type == 'success' &&
    result.data &&
    'result' in result.data &&
    Array.isArray(result.data.result) && // Check if result.data.result is an array
    result.data.result.length > 0 && // Check if the array has elements
    'samples' in result.data &&
    Array.isArray(result.data.samples) &&
    result.data.samples.length > 0 // Check if the array has elements
) {
    captions = JSON.parse(result.data.result[0]);
    samples = result.data.samples;
    console.log(result.data);
}


```


