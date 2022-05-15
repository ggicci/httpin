import { GoPlayProxy } from "@ggicci/goplay"
import CodeBlock from "@theme/CodeBlock"
import React from "react"
import styles from "./GoPlay.module.css"

const GOPLAY_PROXY_URL = "https://goplay.ggicci.me"

type ButtonProps = {
  children: React.ReactNode
  onClick: () => void
}

const Button = (props: ButtonProps) => {
  const { children, onClick } = props
  return (
    <button className={styles.button} onClick={onClick}>
      {children}
    </button>
  )
}

type GoPlayProps = {
  children: React.ReactElement
}

const GoPlay = (props: GoPlayProps) => {
  const { children } = props
  const codeContainer = React.useRef<HTMLDivElement>(null)
  const outputContainer = React.useRef<HTMLDivElement>(null)

  const preEl =
    children && children.props && children.props.mdxType === "pre" && children
  const codeEl = preEl && preEl.props && preEl.props.children

  if (!codeEl || codeEl.props.mdxType !== "code") {
    return <div>GoPlay: the wrapped data is not a codeblock.</div>
  }

  if (!/\blanguage-go\b/.test(codeEl && codeEl.props.className)) {
    return <div>GoPlay: only go code supported.</div>
  }

  function handleRun() {
    if (outputContainer.current) {
      const play = new GoPlayProxy(GOPLAY_PROXY_URL)
      codeContainer.current.classList.remove(styles.hidden)
      play.renderCompile(outputContainer.current, codeEl.props.children.trim())
    }
  }

  async function handleShare() {
    const play = new GoPlayProxy(GOPLAY_PROXY_URL)
    const shareUrl = await play.share(codeEl.props.children.trim())
    window.open(shareUrl, "_blank")
  }

  return (
    <React.Fragment>
      {children}
      <div ref={codeContainer} className={styles.hidden}>
        <CodeBlock language="text">
          <div ref={outputContainer}></div>
        </CodeBlock>
      </div>
      <div className={styles.toolbar}>
        <Button onClick={handleRun}>Run</Button>
        <Button onClick={handleShare}>{"Try it yourself ⇢"}</Button>
      </div>
    </React.Fragment>
  )
}

export default GoPlay
