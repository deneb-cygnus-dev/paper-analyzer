import sys
import logging
import tempfile
import json
from pathlib import Path

from docling_core.types.doc import PictureItem, TableItem, CodeItem

from docling.datamodel.base_models import InputFormat
from docling.datamodel.pipeline_options import PdfPipelineOptions
from docling.document_converter import DocumentConverter, PdfFormatOption

logging.getLogger("docling").setLevel(logging.WARNING)

IMAGE_RESOLUTION_SCALE = 2.0


def main():
    file_path = sys.argv[1]

    input_doc_path = Path(file_path)

    output_dir = Path(tempfile.mkdtemp())

    # Keep page/element images so they can be exported. The `images_scale` controls
    # the rendered image resolution (scale=1 ~ 72 DPI). The `generate_*` toggles
    # decide which elements are enriched with images.
    pipeline_options = PdfPipelineOptions()
    pipeline_options.images_scale = IMAGE_RESOLUTION_SCALE
    pipeline_options.generate_page_images = True
    pipeline_options.generate_picture_images = True

    doc_converter = DocumentConverter(
        format_options={
            InputFormat.PDF: PdfFormatOption(pipeline_options=pipeline_options)
        }
    )
    conv_res = doc_converter.convert(input_doc_path)

    doc_filename = conv_res.input.file.stem

    # Save images of figures and tables
    table_counter = 0
    picture_counter = 0
    code_counter = 0
    metadata = {
        "content": "",
        "tables": [],
        "pictures": [],
        "codes": []
    }
    for element, _level in conv_res.document.iterate_items():
        if isinstance(element, TableItem):
            tag = element.self_ref.replace("/", "@")
            element_image_filename = (
                output_dir / f"{doc_filename}-table-{table_counter}-{tag}.png"
            )
            with element_image_filename.open("wb") as fp:
                element.get_image(conv_res.document).save(fp, "PNG")
            metadata["tables"].append({
                "id": table_counter,
                "path": str(element_image_filename)
            })
            table_counter += 1

        if isinstance(element, PictureItem):
            tag = element.self_ref.replace("/", "@")
            element_image_filename = (
                output_dir / f"{doc_filename}-picture-{picture_counter}-{tag}.png"
            )
            with element_image_filename.open("wb") as fp:
                element.get_image(conv_res.document).save(fp, "PNG")
            metadata["pictures"].append({
                "id": picture_counter,
                "path": str(element_image_filename)
            })
            picture_counter += 1

        if isinstance(element, CodeItem):
            tag = element.self_ref.replace("/", "@")
            element_image_filename = (
                output_dir / f"{doc_filename}-code-{code_counter}-{tag}.png"
            )
            with element_image_filename.open("wb") as fp:
                element.get_image(conv_res.document).save(fp, "PNG")
            metadata["codes"].append({
                "id": code_counter,
                "path": str(element_image_filename)
            })
            code_counter += 1

    content_path = str(output_dir / f"{doc_filename}.json")
    conv_res.document.save_as_json(content_path, None, indent=2)
    metadata["content"] = content_path
    print(json.dumps(metadata, indent=2))


if __name__ == "__main__":
    main()
