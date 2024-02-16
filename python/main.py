import rdflib
import pprint
import argparse
import os
import sys
from fluree.client import FlureeClient

def main():
    print(sys.path)
    parser = argparse.ArgumentParser(description="PyDWIM")
    parser.add_argument('--ingest', type=str, help='Read schemas files in directory')
    args = parser.parse_args()

    if args.ingest:
        print(f'Ingesting from {args.ingest}')
        for file in os.listdir(args.ingest):
            fn, fx = os.path.splitext(os.path.basename(file))

            format = "turtle"
            if file.endswith('.ttl'):
                format = "turtle"
            elif file.endswith('.json'):
                format = "json-ld"
            elif file.endswith('.xsd'):
                format = "application/rdf+xml"
                print(f'Unsupported format for {file}')
                continue
            else:
                print(f'Unknown format for {file}')
                continue

            path = f'{args.ingest}/{file}'
            print(f'Reading {path}')

            g = rdflib.Graph()
            g.parse(f'{args.ingest}/{file}', format=format)
            b = g.serialize(format='json-ld', encoding='utf-8')
            s = b.decode('utf-8')
            
            fc = FlureeClient('http://localhost:8090', f'schemas/{fn}')
            fc.transact(None, s, None)
    else:
        parser.print_help()

if __name__ == "__main__":
    main()
