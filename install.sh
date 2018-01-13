python3.6 -m virtualenv venv
source venv/bin/activate
pip install -r requirements.txt

pushd ui
npm install
popd