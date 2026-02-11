LATEX_BUILDER_IMAGE_NAME ?= azamatbayramov/latex-builder
LATEX_BUILDER_IMAGE_TAG ?= latest
FILE ?= main.tex

.PHONY: pdf clean clean-pdf proposal thesis

# Compile LaTeX file to PDF using the Docker image
pdf:
	docker run --rm -v $(PWD):/workdir -w /workdir $(LATEX_BUILDER_IMAGE_NAME):$(LATEX_BUILDER_IMAGE_TAG) $(FILE)

# Compile LaTeX file to PDF with cleaning before and after
clean-pdf:
	$(MAKE) clean
	$(MAKE) pdf
	$(MAKE) clean

# Remove auxiliary files generated during LaTeX compilation
clean:
	find . -type f \( -name "*.aux" -o -name "*.bbl" -o -name "*.blg" -o \
    -name "*.brf" -o -name "*.fls" -o -name "*.fdb_latexmk" -o \
    -name "*.log" -o -name "*.out" -o -name "*.toc" -o -name "*.lof" -o \
    -name "*.lot" -o -name "*.synctex.gz" -o -name "*.run.xml" \
    -o -name "*.bcf" \) -delete

# Build thesis
thesis:
	$(MAKE) clean-pdf FILE=thesis.tex