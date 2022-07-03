import datetime
import json
import hashlib
import json

import pymongo

class Blockchain:
#Intializing parameters required in the functions
    def __init__(self):
       self.chain = []
       self.prev_hash=0
       self.new_index=0
#creating blocks and adding to the chain
    def create_blockchain(self,timestamp,sender,receiver,value):
        block = {
            'index': self.new_index+ 1,
            'timestamp': timestamp,
            'sender': sender,
            'receiver': receiver,
            'value': value,
            'previous_hash': self.prev_hash
        }
        self.new_index+=1
        self.chain.append(block)
        return block
#Hashing using SHA256 technique
    def hash(self, block):
        encoded_block = json.dumps(block, sort_keys=True).encode()
        return hashlib.sha256(encoded_block).hexdigest()
#displaying the chain
    def get_chain(self):
        print(f"The length of chain is:{len(self.chain)}\nThe chain is {self.chain}")
#saving new blocks and updating database
    def save_db(self):
        myclient = pymongo.MongoClient("mongodb://localhost:27017/")
        mydb = myclient["mydatabase"]
        mycol = mydb["blocks"]
        mycol.drop()
        for i in self.chain:
            mycol.insert_one(i)
#Reading stored blocks from DB
    def read_db(self):
        myclient = pymongo.MongoClient("mongodb://localhost:27017/")
        mydb = myclient["mydatabase"]
        mycol = mydb["blocks"]
        for x in mycol.find():
            print(x)
#Reading transactions from JSON file and adding to block chain
    def read_transaction(self):
        with open('transactions.json') as f:
            data = json.load(f)
        myclient = pymongo.MongoClient("mongodb://localhost:27017/")
        mydb = myclient["mydatabase"]
        mycol = mydb["blocks"]
        #adding pre-present blocks from db
        for x in mycol.find():
            self.new_index=x['index']-1
            self.prev_hash = x['previous_hash']
            block = self.create_blockchain(x['timestamp'], x['sender'], x['receiver'], x['value'])
            self.prev_hash = self.hash(block)
        #adding new data from json to blockchain
        for i in data["transactions"]:
            block = self.create_blockchain(str(datetime.datetime.now()),i['sender'],i['receiver'],i['value'])
            self.prev_hash = self.hash(block)           
#Verifying transaction
    def verify_transaction(self,sender,receiver,value):
        found = False
        myclient = pymongo.MongoClient("mongodb://localhost:27017/")
        mydb = myclient["mydatabase"]
        mycol = mydb["blocks"]
        for x in mycol.find():
            if(x['sender']==sender and x['receiver']==receiver and x['value']==value):
                found = True
                break
        return found


#MAIN FUNCTION
#To check selected option
ch=0
bc = Blockchain()
while ch!=5:
    ch = int(input("Enter 1.To read from JSON\n2.Store to DB\n3.Read from DB\n4.Check validity of transaction\n5.Exit\n")) 
    if(ch==1):
        bc.read_transaction()
        bc.get_chain()
    if(ch==2):
        bc.save_db()
    if(ch==3):
        bc.read_db()
    if(ch==4):
        sender = input("Enter Sender's Name: ")
        receiver = input("Enter Receiver's name: ")
        value = input("Enter amount sent: ")
        if(bc.verify_transaction(sender,receiver,value)):
            print("The transaction is valid")
        else:
            print("Transaction Invalid")


