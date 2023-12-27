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
    Directive: "required",
    Field: "access_token",
    Key: "",
    Value: nil,
    ErrorMessage: "missing required field",
}
`,
      },
    ],
  },
}
