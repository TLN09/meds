from Crypto.Hash import SHAKE256
from sys import argv

shake = SHAKE256.new()
shake.update(argv[1].encode())

for _ in range(int(argv[2])):
    print(ord(shake.read(1)))
