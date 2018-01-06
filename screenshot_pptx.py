import os.path
import subprocess
import sys
import xml.etree.ElementTree as ET
import zipfile

import database

def copy_in_zip(zip_handle, from_path, to_path):
    with zip_handle.open(from_path) as read_slide:
        to_copy = read_slide.read()
    with zip_handle.open(to_path, 'w') as write_slide:
        write_slide.write(to_copy)
    return to_copy

def generate_thumbs(powerpoint_file):
    original_filename = os.path.realpath(powerpoint_file)
    
    session = database.Session()
    doc = database.Document()
    doc.path = original_filename
    session.add(doc)

    subprocess.run(['cp', powerpoint_file, 'trashed.pptx'])
    powerpoint_file = 'trashed.pptx'
    with zipfile.ZipFile(powerpoint_file) as pptx_zip:
        slide_count = len([i.filename for i in pptx_zip.infolist() if i.filename.startswith('ppt/slides/') and not i.filename.startswith('ppt/slides/_rels')])

    print(slide_count)

    for slide_number in range(1, slide_count+1):
        path = 'ppt/slides/slide' + str(slide_number) + '.xml'
        rels_path = 'ppt/slides/_rels/slide' + str(slide_number) + '.xml.rels'

        print(path)

        with zipfile.ZipFile(powerpoint_file, 'a') as pptx_zip:
            print("slides:")
            print([i.filename for i in pptx_zip.infolist() if i.filename.startswith('ppt/slides/') and not i.filename.startswith('ppt/slides/_rels')])
            if slide_number > 1:
                xml_content = copy_in_zip(pptx_zip, path, 'ppt/slides/slide1.xml')
                copy_in_zip(pptx_zip, rels_path, 'ppt/slides/_rels/slide1.xml.rels')
            else:
                with pptx_zip.open(path) as read_slide:
                    xml_content = read_slide.read()

        thumb_name = original_filename + '_slide' + str(slide_number) + '.png'
        subprocess.run(['qlmanage', '-t', '-s', '400', '-o', '.', powerpoint_file])
        
        thumbnail = database.Slide()
        thumbnail.document = doc
        thumbnail.slide = slide_number
        
        text_content = scrape_text_content(xml_content)

        slide_content = database.SlideContent()
        slide_content.content = text_content
        slide_content.slide = thumbnail
        session.add(slide_content)


        with open('trashed.pptx.png', 'rb') as png_file:
            thumbnail.thumnail_png = png_file.read()
        session.add(thumbnail)

    session.commit()

def scrape_text_content(xml_content):
    tree = ET.fromstring(xml_content)

    return '\n'.join([text_elem.text for text_elem in tree.iter('{http://schemas.openxmlformats.org/drawingml/2006/main}t')])

if __name__ == '__main__':
    filename = sys.argv[1]
    generate_thumbs(filename)

