import zipfile
import subprocess

def copy_in_zip(zip_handle, from_path, to_path):
    with zip_handle.open(from_path) as read_slide:
        to_copy = read_slide.read()
    with zip_handle.open(to_path, 'w') as write_slide:
        write_slide.write(to_copy)

def generate_thumbs(powerpoint_file):
    original_filename = powerpoint_file
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
                copy_in_zip(pptx_zip, path, 'ppt/slides/slide1.xml')
                copy_in_zip(pptx_zip, rels_path, 'ppt/slides/_rels/slide1.xml.rels')

        thumb_name = original_filename + '_slide' + str(slide_number) + '.png'
        subprocess.run(['qlmanage', '-t', '-s', '400', '-o', thumb_name, powerpoint_file])
        subprocess.run(['cp', 'trashed.pptx.png', thumb_name])

if __name__ == '__main__':
    generate_thumbs('some_presentation.pptx')
