import rdflib
import pprint
import fluree

def main():
    print("Hello World!")
    g = rdflib.Graph()
    g.parse("schemas/usefulinc/doap.xml")
    g.parse("schemas/w3_org/owl.ttl")
    g.parse("schemas/w3_org/prov.ttl")
    g.parse("schemas/w3_org/rdf.ttl")
    g.parse("schemas/w3_org/rdfs.ttl")
    g.parse("schemas/w3_org/skos.ttl")
    #g.parse("schemas/w3_org/xsd.xsd", format="xml")
    ld = g.serialize(format="json-ld")
    #pprint.pprint(ld)

    fluree.schema_gen(ld, "gen/schemas/ld-json/")

if __name__ == "__main__":
    main()
