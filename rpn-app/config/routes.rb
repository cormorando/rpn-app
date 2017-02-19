Rails.application.routes.draw do
  get 'main/index'
  post 'main/parse'

  root 'main#index'
end
