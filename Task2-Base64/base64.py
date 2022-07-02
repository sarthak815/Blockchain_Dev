#Base64 Index table
base_64_id = ['A', 'B', 'C', 'D', 'E', 'F', 'G', 'H',
                                'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P',
                                'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X',
                                'Y', 'Z', 'a', 'b', 'c', 'd', 'e', 'f',
                                'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n',
                                'o', 'p', 'q', 'r', 's', 't', 'u', 'v',
                                'w', 'x', 'y', 'z', '0', '1', '2', '3',
                                '4', '5', '6', '7', '8', '9', '+', '/']
#Obtain list of binary values
def ascii_list(ip):
    ascii_val  = list()
    for a in ip:
        ascii_val.append(ord(a))

    return ascii_val
#ascii to binary string
def ascii_bin(ascii_list):
    op = ""
    for i in ascii_list:
        bindata = f'{i:08b}'
        op+=bindata
    return op
#binary to serialized string
def bin_serialized(bin_string,n):
    split_arr = [bin_string[i:i + 6] for i in range(0, len(bin_string), 6)]
    op = ""
    int_arr = list()
    for ch in split_arr:
        if(len(ch)!=6):
            while(len(ch)!=6):
                ch = ch+'0'
        x = int(ch, 2)
        int_arr.append(x)
    for val in int_arr:
        op+=base_64_id[int(val)]
    if n%3!=0:
        op+='='
    return op
#Obtaining output
def exec(ip):
    for a in ip:
        if((a>='a'and a<='z') or (a>='A'and a<='Z')or a=='+' or a=='/' or (a>='0' and a<='9')):
            continue
        else:
            print("Invalid Character: Cannot be serialized")
            exit(0)
            break
    ascii_val = ascii_list(ip)
    bin_val = ascii_bin(ascii_val)
    return bin_serialized(bin_val,len(ip))
#input to be serialized
ip = input("Enter string to be serialized: ")
#Ensuring input can be serialized
print(exec(ip))

