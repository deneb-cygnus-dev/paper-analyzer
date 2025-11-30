from docling.document_converter import DocumentConverter

source = "testdata/artifacts/Constrained Detecting Arrays.pdf"
converter = DocumentConverter()
doc = converter.convert(source).document

print(doc.export_to_markdown())
