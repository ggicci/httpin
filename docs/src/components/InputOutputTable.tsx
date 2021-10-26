import CodeBlock from "@theme/CodeBlock"
import React from "react"

interface Row {
  input: string
  output: string
}

interface Props {
  inputTitle: string
  outputTitle: string
  rows: Row[]
}

const Index = (props: Props) => {
  return (
    <table>
      <thead>
        <tr>
          <th>{props.inputTitle || "Input"}</th>
          <th>{props.outputTitle || "Output"}</th>
        </tr>
      </thead>
      <tbody>
        {props.rows.map((row, i) => (
          <tr key={i}>
            <td>
              <CodeBlock>{row.input.trim()}</CodeBlock>
            </td>
            <td>
              <CodeBlock className="language-go">{row.output.trim()}</CodeBlock>
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  )
}

export default Index
