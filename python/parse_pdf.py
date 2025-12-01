import sys
import tempfile
import logging

from docling_core.types.doc.document import SectionHeaderItem, DoclingDocument
from docling.datamodel.pipeline_options import PdfPipelineOptions
from docling.document_converter import DocumentConverter, PdfFormatOption
from docling.datamodel.base_models import InputFormat
from docling_core.types.doc.labels import DocItemLabel
from docling_core.types.doc import ImageRefMode


# Suppress docling logs
logging.getLogger("docling").setLevel(logging.WARNING)

IMAGE_RESOLUTION_SCALE = 1.0


def setup_converter():
    pipeline_options = PdfPipelineOptions()
    pipeline_options.images_scale = IMAGE_RESOLUTION_SCALE
    pipeline_options.generate_page_images = False
    pipeline_options.generate_picture_images = True

    doc_converter = DocumentConverter(
        format_options={
            InputFormat.PDF: PdfFormatOption(pipeline_options=pipeline_options)
        }
    )

    return doc_converter


def determine_paper_template(document: DoclingDocument):
    for element, _level in document.iterate_items():
        if element.label != DocItemLabel.SECTION_HEADER:
            continue
        if all(x in element.text.lower() for x in ["acm", "reference", "format"]):
            return "acm_conference_template"
    return "unknown"


def clear_redundant_elements(document: DoclingDocument):
    elements_to_delete = []
    for element, _level in document.iterate_items():
        if element.label == DocItemLabel.FOOTNOTE:
            elements_to_delete.append(element)
        if element.label == DocItemLabel.PAGE_HEADER:
            elements_to_delete.append(element)
        if element.label == DocItemLabel.PAGE_FOOTER:
            elements_to_delete.append(element)
    document.delete_items(node_items=elements_to_delete)
    return document


def get_section_headers_acm_conference_template(document: DoclingDocument):
    section_headers = []
    for element, _level in document.iterate_items():
        if element.label == DocItemLabel.SECTION_HEADER:
            section_headers.append(element)
    return section_headers


def process_acm_conference_template(document: DoclingDocument):
    # get the title
    section_headers = get_section_headers_acm_conference_template(document)
    title = section_headers[0].text
    document = clear_redundant_elements(document)
    return document


def main():
    pdf_path = sys.argv[1]

    output_dir = tempfile.mkdtemp()

    doc_converter = setup_converter()
    conv_res = doc_converter.convert(pdf_path)

    document = conv_res.document

    paper_template = determine_paper_template(document)

    if paper_template == "acm_conference_template":
        document = process_acm_conference_template(document)

    with open("output-wo-footnotes.md", "w") as f:
        f.write(document.export_to_markdown())
    with open("output.txt", "w") as f:
        for element, _level in document.iterate_items():
            if isinstance(element, SectionHeaderItem):
                # section headers
                pass
            f.write(str(element) + "\n")
    document.save_as_json(
        "output.json", None, image_mode=ImageRefMode.REFERENCED, indent=2
    )
    print(output_dir)


if __name__ == "__main__":
    main()
