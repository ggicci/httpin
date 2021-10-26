export default {
  inputOutput: {
    inputTitle: "Request (URL query)",
    outputTitle: "ListUsersInput",
    rows: [
      {
        input: `
GET /users?is_member=1&age_range[]=3&age_range[]=5
`,
        output: `
{
    IsMember: true,
    AgeRange: []int{3, 5},
}`,
      },
      {
        input: `
GET /users?age_range=3&age_range=5
`,
        output: `
{
    IsMember: false,
    AgeRange: []int{3, 5},
}`,
      },
      {
        input: `
GET /users?is_member=true
`,
        output: `
{
    IsMember: true,
    AgeRange: []int{},
}
`,
      },
    ],
  },
}
