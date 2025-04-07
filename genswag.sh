#!/bin/bash

# Fungsi untuk mendapatkan input dari pengguna
get_input() {
  read -p "$1: " input
  echo "$input"
}

# Fungsi untuk menghasilkan komentar Swag
generate_swag_comment() {
  summary=$(get_input "Summary")
  description=$(get_input "Description")
  tags=$(get_input "Tags (comma-separated)")
  accept=$(get_input "Accept (e.g., json)")
  produce=$(get_input "Produce (e.g., json)")
  security=$(get_input "Security (e.g., BearerAuth)")

  echo "// @Summary $summary"
  echo "// @Description $description"
  echo "// @Tags $tags"
  echo "// @Accept $accept"
  echo "// @Produce $produce"
  echo "// @Security $security"

  # Parameter query dinamis
  while true; do
    param_name=$(get_input "Parameter query name (or press Enter to finish)")
    if [ -z "$param_name" ]; then
      break
    fi
    param_type=$(get_input "Parameter query type (e.g., int, string, object)")
    param_required=$(get_input "Parameter query required (true/false)")
    param_description=$(get_input "Parameter query description")
    param_default=$(get_input "Parameter query default (optional)")

    echo "// @Param $param_name query $param_type $param_required \"$param_description\" $([ -n "$param_default" ] && echo "default($param_default)")"
  done

  # Response success dinamis
  response_code=$(get_input "Response success code (e.g., 200)")
  response_object=$(get_input "Response success object (e.g., response.BaseResponse{data=database.Paginator{records=school.TypeExam}})")
  response_description=$(get_input "Response success description")

  echo "// @Success $response_code {object} $response_object \"$response_description\""

  # Router dinamis
  router=$(get_input "Router path (e.g., /academic/exam/type-exam)")
  http_method=$(get_input "HTTP method (e.g., get, post, put, delete)")

  echo "// @Router $router [$http_method]"
}

# Memanggil fungsi untuk menghasilkan komentar Swag
generate_swag_comment