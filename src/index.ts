import PromptSync from "prompt-sync"
import { program } from "commander"
import path from "path"
import fs from "fs"

const inputUtil = PromptSync({ sigint: true })

const variables: { [key: string]: any } = {}
const statements = ["ㅇㅉ", "ㅈㅉ", "어쩔", "티비"]

const execute = async (code: string) => {
  const lines: string[] = code.split("\n")
  if (lines[0] !== "어쩔티비") throw new Error("아무것도 모르죠?")
  if (lines[lines.length - 1] !== "저쩔티비") throw new Error("아무것도 모르죠?")
  for (const line of lines) {
    const [statement] = line.split(" ")
    if (statement === "ㅇㅉ") {
      allocateVariable(line)
    } else if (statement === "어쩔") {
      print(line)
    } else if (statement === "티비") {
      input(line)
    } else if (statement === "ㅇㅇ") {
      assignVariable(line)
    }
  }
}

const getVariable = (line: string) => {
  // 키워드는 변수이름이 될 수 없습니다
  if (line.split(" ")[0] === "티비") return input(line)
  if (statements.includes(line.split(" ")[0])) return
  const value = variables[line]
  if (!value) return null
  return value
}

const allocateVariable = (line: string) => {
  const [statement, name, first, ...values] = line.split(" ")
  if (statement !== "ㅇㅉ") return
  if (!name || name.length <= 0) throw new Error("어쩔변수")
  let allocatingValue = ""
  if (first === "티비") {
    const inputValue = input([first, ...values].join(" "))
    if (inputValue) allocatingValue = inputValue
  } else {
    allocatingValue = String(toNumber(first))
  }
  variables[name] = allocatingValue
}

const assignVariable = (line: string) => {
  const [statement, name, ...values] = line.split(" ")
  if (statement !== "ㅈㅉ") return
  if (!name || name.length <= 0) throw new Error("어쩔변수")
  const doesVariableExist = getVariable(name)
  if (doesVariableExist === null) throw new Error("어쩔변수")
  let value = ""
  if (name === "티비") {
    const inputValue = input(values.join(" "))
    if (inputValue) value = inputValue
  } else {
    value = String(toNumber(values.join(" ")))
  }
  variables[name] = value
}

const toNumber = (line: string) => {
  const pluses = line.split("").filter((v) => v === "어").length
  const minuses = line.split("").filter((v) => v === "쩔").length
  return pluses - minuses
}

const print = (line: string) => {
  const [statement, ...printOut] = line.split(" ")
  if (statement !== "어쩔") return
  console.log(printOut.map((v) => getVariable(v)).join(" "))
}

const input = (line: string) => {
  const [statement, ...inputString] = line.split(" ")
  if (statement !== "티비") return
  const inputUser = inputUtil(inputString.join(" ") + "\n")
  return inputUser
}

program.parse()

const targetFilePath = path.join(process.cwd(), program.args.join(" "))
if (!fs.existsSync(targetFilePath)) {
  throw new Error("어쩔파일")
}
const codes = fs.readFileSync(targetFilePath)
execute(codes.toString("utf-8"))
