export default {
  inputOutput: {
    inputTitle: "Request",
    outputTitle: "ListTasksQuery",
    rows: [
      {
        input: `
  GET /tasks?page=4&perPage=10&state=failed&state=succeeded
  `,
        output: `
  {
      Page:      4,
      PerPage:   10,
      StateList: []string{"failed", "succeeded"},
  }`,
      },
      {
        input: `
  GET /tasks
  `,
        output: `
  {
            Page:      1,
            PerPage:   20,
            StateList: []string{"pending", "running"},
  }`,
      },
    ],
  },
}
