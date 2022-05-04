import random
from worker import app


@app.task(name='fibonacci')
def task_fib(n: int):
    return fibonacci(n)


def fibonacci(n: int) -> int:
    print(f"Calculating Fibonacci of {n}")
    if n == 0:
        return 0
    elif n == 1:
        return 1
    else:
        a = 0
        b = 1

        for i in range(1, n):
            c = a + b
            a = b
            b = c
        return b
