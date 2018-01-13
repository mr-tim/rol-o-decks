from sqlalchemy import Column, Binary, DateTime, Integer, String, ForeignKey, Text, create_engine
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import sessionmaker, relationship

Base = declarative_base()

class Document(Base):
    __tablename__ = 'document'

    id = Column(Integer, primary_key=True)
    path = Column(String)
    created = Column(DateTime)
    last_modified = Column(DateTime)

    thumbnails = relationship("Slide", order_by="Slide.slide", back_populates='document', cascade='all')


class Slide(Base):
    __tablename__ = 'slide'

    id = Column(Integer, primary_key=True)
    thumnail_png = Column(Binary)
    slide = Column(Integer)
    document_id = Column(Integer, ForeignKey('document.id'))

    document = relationship("Document", back_populates="thumbnails")
    content = relationship("SlideContent", back_populates="slide", cascade='all')

class SlideContent(Base):
    __tablename__ = 'slide_content'

    slide_id = Column(Integer, ForeignKey('slide.id'), primary_key=True)
    content = Column(Text)

    slide = relationship("Slide", back_populates="content")

engine = create_engine('sqlite:///db/rolodecks.sqlite', echo=True)

engine.execute('create virtual table if not exists slide_content using fts4(slide_id INTEGER, content TEXT);')

Base.metadata.create_all(engine)

Session = sessionmaker(bind=engine)