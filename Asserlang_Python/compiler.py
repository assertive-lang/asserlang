import exceptions
from enum import Enum, auto


class KeyWordType(Enum):
    BLANK = auto()
    VAR_DECLARE = auto()
    VAR_ASSIGN = auto()
    ASCII_VAR_DECLARE = auto()
    ASCII_VAR_ASSIGN = auto()
    INPUT = auto()
    OUTPUT = auto()
    IF = auto()
    FUNCTION_DECLARE = auto()
    FUNCTION_CALL = auto()
    FUNCTION_RETURN = auto()
    GOTO = auto()


class AsserLang:
    def __init__(self):
        self.index = 0
        self.int_data = dict()
        self.string_data = dict()

    def get_int(self, code):
        plus = code.count('ㅋ')
        minus = code.count('ㅎ')

        code = code.replace('ㅋ', '')
        code = code.replace('ㅎ', '')

        if not code.empty():
            raise exceptions.VariableException()

        return plus - minus

    def get_number(self, code):
        output = 0
        now_add = True

        while len(code) != 0:
            if code.startswith('ㅌㅂ'):  # input
                if now_add:
                    output += int(input())
                else:
                    output *= int(input())
                    now_add = True

                code = code[2:]
            elif code.startswith('ㅋ') or code.startswith('ㅎ'):
                now_char = code[0]
                now_num = 0

                while now_char == 'ㅋ' or now_char == 'ㅎ':
                    if not code:
                        break

                    now_char = code[0]

                    if now_char == 'ㅋ':
                        now_num += 1
                    elif now_char == 'ㅎ':
                        now_num -= 1
                    else:
                        break

                    code = code[1:]

                if now_add:
                    output += now_num
                else:
                    output *= now_num
                    now_add = True
            elif code.startswith('ㅌ'):
                now_add = False
                code = code[1:]

            else:
                is_var_name = False

                for key in self.int_data.keys():
                    if code.startswith(key):
                        if now_add:
                            output += self.int_data[key]
                        else:
                            output *= self.int_data[key]
                            now_add = True

                        code = code[len(key):]
                        is_var_name = True

                if not is_var_name:
                    raise exceptions.VariableException()

        return output

    def get_output(self, code):
        for keys in self.string_data.keys():
            if code.startswith(keys):
                if code == keys:
                    return self.string_data[keys]
                else:
                    raise exceptions.VariableException()

        return self.get_number(code)

    @staticmethod
    def get_type(self, code):
        if not code:
            return KeyWordType.BLANK
        elif code.startswith('어쩔'):  # 변수 선언
            return KeyWordType.VAR_DECLARE
        elif code.startswith('저쩔'):  # 변수 할당
            return KeyWordType.VAR_ASSIGN
        elif code.startswith('우짤래미'):  # 아스키 변수 선언
            return KeyWordType.ASCII_VAR_DECLARE
        elif code.startswith('저짤래미'):  # 아스키 변수 할당
            return KeyWordType.ASCII_VAR_ASSIGN
        elif code.startswith('ㅌㅂ'):  # input
            return KeyWordType.INPUT
        elif code.startswith('ㅇㅉ'):  # 출력
            return KeyWordType.OUTPUT
        elif code.startswith('안물'):  # 함수 선언
            return KeyWordType.FUNCTION_DECLARE
        elif code.startswith('안궁'):  # 함수 호출
            return KeyWordType.FUNCTION_CALL
        elif code.startswith('무지개반사'):  # 함수 반환
            return KeyWordType.FUNCTION_RETURN
        elif code.startswith('화났쥬?'):  # 조건문
            return KeyWordType.IF
        elif code.startswith(';;'):  # GOTO
            return KeyWordType.GOTO
        else:
            raise SyntaxError()

    def compileLine(self, code):
        if code == '':
            return None

        TYPE = self.get_type(self, code)

        if TYPE == KeyWordType.BLANK:
            return

        if TYPE == KeyWordType.VAR_DECLARE:
            # 어쩔{변수명}~{변수대입값}

            code = code.removeprefix('어쩔')
            components = code.split('~')
            length = len(components)

            val = 0
            varname = components[0]

            if length != 1 and length != 2:
                raise exceptions.VariableException()

            if varname in self.int_data:
                raise exceptions.VariableException()

            if length == 2:
                try:
                    val = self.get_number(components[1])
                except Exception:
                    raise exceptions.VariableException()

            self.int_data[varname] = val
            # print(f'declared val on : ({varname}, {val})')

        elif TYPE == KeyWordType.VAR_ASSIGN:
            # 저쩔{변수명}~{변수대입값}

            code = code.removeprefix('저쩔')
            components = code.split('~')
            length = len(components)

            val = 0
            varname = components[0]

            if length != 1 and length != 2:
                raise exceptions.VariableException()

            if length == 2:
                try:
                    val = self.get_number(components[1])
                except Exception:
                    raise exceptions.VariableException()

            if not (varname in self.int_data):
                raise exceptions.VariableException()

            self.int_data[varname] = val
            # print(f'assigned val on : ({varname}, {val})')

        elif TYPE == KeyWordType.ASCII_VAR_DECLARE:
            # 우짤래미{변수명}~{변수대입값}

            code = code.removeprefix('우짤래미')
            components = code.split('~')
            length = len(components)

            val = 0
            varname = components[0]

            if length != 1 and length != 2:
                raise exceptions.VariableException()

            if length == 2:
                try:
                    val = self.get_number(components[1])
                except Exception:
                    raise exceptions.VariableException()

            if varname in self.string_data:
                raise exceptions.VariableException()

            self.string_data[varname] = chr(val)
            # print(f'declared ascii_val on : ({varname}, {val})')

        elif TYPE == KeyWordType.ASCII_VAR_ASSIGN:
            # 저짤래미{변수명}~{변수대입값}

            code = code.removeprefix('저짤래미')
            components = code.split('~')
            length = len(components)

            val = 0
            varname = components[0]

            if length != 1 and length != 2:
                raise exceptions.VariableException()

            if length == 2:
                try:
                    val = self.get_number(components[1])
                except Exception:
                    raise exceptions.VariableException()

            if not (varname in self.string_data):
                raise exceptions.VariableException()

            self.string_data[varname] = chr(val)
            # print(f'assigned ascii_val on : ({varname}, {val})')

        elif TYPE == KeyWordType.INPUT:
            input()
        elif TYPE == KeyWordType.OUTPUT:
            # ㅇㅉ{출력값}

            code = code.removeprefix('ㅇㅉ')
            try:
                val = self.get_output(code)
                print(val, end='')
            except Exception:
                raise exceptions.VariableException()

        elif TYPE == KeyWordType.GOTO:
            number = code[2:]
            number_int = self.get_number(number)

            return number_int - 2

        elif TYPE == KeyWordType.IF:
            code = code.removeprefix('화났쥬?')
            components = code.split('킹받쥬?')

            if len(components) != 2:
                raise exceptions.ConStatementException()

            num = components[0]
            line = components[1]

            num = self.get_number(num)

            if num == 0:
                states = self.compileLine(line)

                if type(states) == int:
                    self.index = states - 1 # index + 1 on main compiler

    def compile(self, code):
        self.index = 0

        if code.pop(0) != '쿠쿠루삥뽕' or code.pop() != '슉슈슉슉':
            raise exceptions.MainException()

        while self.index < len(code):
            # print(f'compiling line {self.index + 2}... {code[self.index]}')

            states = self.compileLine(code[self.index])
            self.index += 1

            if type(states) == int:
                self.index = states

    def compile_file(self, path):
        try:
            with open(path, encoding='utf-8-sig') as mte_file:
                code = mte_file.read().splitlines()
                self.compile(code)

        except FileNotFoundError:
            raise exceptions.FileNotFoundException()

        except UnicodeDecodeError:
            raise exceptions.FileNotFoundException()

if __name__ == '__main__':
    compiler = AsserLang()
    compiler.compile_file('test.astv')
