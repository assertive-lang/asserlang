class MainException(Exception):
    def __init__(self):
        super().__init__('아무것도 모르죠?')


class VariableException(Exception):
    def __init__(self):
        super().__init__('어쩔변수')


class FileNotFoundException(Exception):
    def __init__(self):
        super().__init__('어쩔파일')


class FunctionException(Exception):
    def __init__(self):
        super().__init__('안물')


class ConStatementException(Exception):
    def __init__(self):
        super().__init__('어쩔조건')


class GotoException(Exception):
    def __init__(self):
        super().__init__('어쩔GOTO인덱스')


class VariableNamingException(Exception):
    def __init__(self):
        super().__init__('어쩔변수이름')
