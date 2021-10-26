export default {
  inputOutput: {
    inputTitle: "Request (Mixed, Body + URL query)",
    outputTitle: "Profile",
    rows: [
      {
        input: `
GET /users?role=backend&hireable=true
`,
        output: `
{
  Role: "backend",
  Hireable: true,
}

// works just like query directive
`,
      },
      {
        input: `
POST /users HTTP/1.1
Host: foo.example
Content-Type: application/x-www-form-urlencoded

role=frontend&hireable=false
`,
        output: `
{
  Role: "frontend",
  Hireable: false,
}
`,
      },
      {
        input: `
POST /users?hireable=true HTTP/1.1
Host: foo.example
Content-Type: application/x-www-form-urlencoded

role=frontend&hireable=false
`,
        output: `
{
  Role: "frontend",
  Hireable: false,
}

// Hireable in the body overrides the URL query.
`,
      },
      {
        input: `
POST /users?hireable=true HTTP/1.1
Host: foo.example
Content-Type: application/x-www-form-urlencoded

role=frontend
`,
        output: `
{
  Role: "frontend",
  Hireable: true,
}

// Hireable is from URL query.
`,
      },
    ],
  },
}
