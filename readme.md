


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

```js

let params: Params = {
    tone: "",
    words: "",
    hashtags: true,
    emojis: true,
    network: "",
    context: "",
    posts: "",
    url: ""
}

let meeting: Meeting = {
    startTime: "",
    endTime: "",
    link: ""
}

let prompt = "*The user prompt goes here*"

const data = JSON.stringify({
    prompt: prompt,
    params: params,
    meeting: meeting
})

const res = await fetch("http://127.0.0.1:8080/text", {
    method: "POST",
    body: data,
})

```

