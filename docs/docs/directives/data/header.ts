export default {
  inputOutput: {
    inputTitle: "Request (Header)",
    outputTitle: "TokenInput",
    rows: [
      {
        input: `
GET /users HTTP/1.1

Host: foo.example
X-API-Token: abc
`,
        output: `{ Token: "abc" }`,
      },
      {
        input: `
GET /users HTTP/1.1

Host: foo.example
x-api-token: abc
`,
        output: `{ Token: "abc" }`,
      },
      {
        input: `
GET /users HTTP/1.1

Host: foo.example
Authorization: good
`,
        output: `{ Token: "good" }`,
      },
      {
        input: `
GET /users HTTP/1.1

Host: foo.example
X-Api-Token: apple
Authorization: banana
`,
        output: `
{ Token: "apple" }

// Check order: X-Api-Token -> Authorization
`,
      },
    ],
  },
}
