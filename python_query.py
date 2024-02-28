from pymongo import MongoClient

def query_blob_hash(collection_name, blob_hash):
    # Connect to MongoDB
    client = MongoClient('localhost', 27017)

    # Access the database
    db = client['blob_storage']

    # Access the collection
    collection = db[collection_name]

    # Query the collection for the blobHash
    result = collection.find_one({'blob_hash': blob_hash})

    # Close the connection
    client.close()

    return result

def main():
    collection_name = 'blobs'
    blob_hash = '0x016dfd20188f4486eda9f09ee8917727b2e42bdc50e02f2d115120e323627e69'

    result = query_blob_hash(collection_name, blob_hash)
    if result:
        print('Found result:', result)
    else:
        print('Blob hash not found')

if __name__ == '__main__':
    main()
