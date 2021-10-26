export default {
  inputOutput: {
    inputTitle: "Request (Body)",
    outputTitle: "CreateUserInput",
    rows: [
      {
        input: `
POST /users HTTP/1.1
Host: foo.example
Content-Type: application/json

{ "login": "alex", "gender": "female" }
  `,
        output: `
Payload: &User{
    Login: "alex",
    Gender: "female",
}
`,
      },
    ],
  },
}
