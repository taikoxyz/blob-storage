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
    blob_hash = 'YOUR_BLOB_HASH_TO_QUERY'

    result = query_blob_hash(collection_name, blob_hash)
    if result:
        print('Found result:', result)
    else:
        print('Blob hash not found')

if __name__ == '__main__':
    main()
