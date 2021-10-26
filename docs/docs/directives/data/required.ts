export default {
  inputOutput: {
    inputTitle: "Request",
    outputTitle: "Output",
    rows: [
      {
        input: `
GET /users?access_token=abc&page=1
  `,
        output: `
&TokenInput{
    Token: "abc",
}
`,
      },
      {
        input: `
GET /users?page=1
  `,
        output: `
// error occurred
&InvalidFieldError{
    Field: "access_token",
    Source: "required",
    Value: nil,
    ErrorMessage: "missing required field",
}
`,
      },
    ],
  },
}
