
truncate usuario CASCADE;

INSERT INTO public.usuario (email, nome, senha_hash) VALUES ('teste@email.com', 'teste', '$2a$14$EbWwzxcBWRXxpxOFb7vEUe3b1Z.sp7Gn5A.GfFDE4F/41MISwr7Oy' /* teste */);
