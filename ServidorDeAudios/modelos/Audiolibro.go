package modelos

type Audiolibro struct {
	idAudio   int
	titulo    string
	autor     string
	narrador  string
	editorial string
	isbn      string
	capitulo  string
}

func (a *Audiolibro) GetIdAudio() int              { return a.idAudio }
func (a *Audiolibro) SetIdAudio(id int)             { a.idAudio = id }

func (a *Audiolibro) GetTitulo() string             { return a.titulo }
func (a *Audiolibro) SetTitulo(titulo string)        { a.titulo = titulo }

func (a *Audiolibro) GetAutor() string              { return a.autor }
func (a *Audiolibro) SetAutor(autor string)          { a.autor = autor }

func (a *Audiolibro) GetNarrador() string           { return a.narrador }
func (a *Audiolibro) SetNarrador(narrador string)    { a.narrador = narrador }

func (a *Audiolibro) GetEditorial() string          { return a.editorial }
func (a *Audiolibro) SetEditorial(editorial string)  { a.editorial = editorial }

func (a *Audiolibro) GetIsbn() string               { return a.isbn }
func (a *Audiolibro) SetIsbn(isbn string)            { a.isbn = isbn }

func (a *Audiolibro) GetCapitulo() string           { return a.capitulo }
func (a *Audiolibro) SetCapitulo(capitulo string)    { a.capitulo = capitulo }
