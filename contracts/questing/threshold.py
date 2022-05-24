"""
Thesis: yield the appropriate threshold indice from a given input percentile.
n = 30

x | 0 -  1 -  2 -  3 -   4
--------------------------
y | 0 - 20 - 50 - 15 -  15
--------------------------
z | 0 - 20 - 70 - 85 - 100
"""

def convict(n, percentiles):
    acc = 0
    for i, pct in enumerate(percentiles):
        prev_acc = acc
        acc += pct
        if prev_acc <= n <= acc:
            return i

if __name__ == "__main__":
    n = 45

    percentiles = [ 10, 10, 50, 15, 15]
    indice = convict(n, percentiles)
    print(indice, percentiles[indice])
