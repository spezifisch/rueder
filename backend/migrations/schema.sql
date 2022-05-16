--
-- PostgreSQL database dump
--

-- Dumped from database version 13.6 (Debian 13.6-1.pgdg110+1)
-- Dumped by pg_dump version 13.6

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: articles; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.articles (
    id uuid NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    seq bigint NOT NULL,
    feed_id uuid NOT NULL,
    site_guid character varying(2048) NOT NULL,
    posted_at timestamp without time zone NOT NULL,
    link character varying(2048),
    thumbnail character varying(2048),
    image character varying(2048),
    image_title character varying(2048),
    title text,
    teaser text,
    content jsonb NOT NULL
);


ALTER TABLE public.articles OWNER TO postgres;

--
-- Name: articles_seq_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.articles_seq_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.articles_seq_seq OWNER TO postgres;

--
-- Name: articles_seq_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.articles_seq_seq OWNED BY public.articles.seq;


--
-- Name: feeds; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.feeds (
    id uuid NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    fetched_at timestamp without time zone NOT NULL,
    fetch_delay_s integer NOT NULL,
    fetcher_state jsonb NOT NULL,
    feed_url character varying(2048) NOT NULL,
    site_url character varying(2048),
    title character varying(1024),
    icon character varying(2048)
);


ALTER TABLE public.feeds OWNER TO postgres;

--
-- Name: labels; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.labels (
    id uuid NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    title character varying(1024),
    color character varying(16)
);


ALTER TABLE public.labels OWNER TO postgres;

--
-- Name: schema_migration; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.schema_migration (
    version character varying(14) NOT NULL
);


ALTER TABLE public.schema_migration OWNER TO postgres;

--
-- Name: user_feeds; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.user_feeds (
    id integer NOT NULL,
    user_id uuid NOT NULL,
    feed_id uuid NOT NULL
);


ALTER TABLE public.user_feeds OWNER TO postgres;

--
-- Name: user_feeds_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.user_feeds_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.user_feeds_id_seq OWNER TO postgres;

--
-- Name: user_feeds_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.user_feeds_id_seq OWNED BY public.user_feeds.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id uuid NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    auth_origin character varying(255) NOT NULL,
    auth_subject character varying(255) NOT NULL,
    folders jsonb NOT NULL
);


ALTER TABLE public.users OWNER TO postgres;

--
-- Name: articles seq; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.articles ALTER COLUMN seq SET DEFAULT nextval('public.articles_seq_seq'::regclass);


--
-- Name: user_feeds id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_feeds ALTER COLUMN id SET DEFAULT nextval('public.user_feeds_id_seq'::regclass);


--
-- Name: articles articles_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.articles
    ADD CONSTRAINT articles_pkey PRIMARY KEY (id);


--
-- Name: feeds feeds_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.feeds
    ADD CONSTRAINT feeds_pkey PRIMARY KEY (id);


--
-- Name: labels labels_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.labels
    ADD CONSTRAINT labels_pkey PRIMARY KEY (id);


--
-- Name: user_feeds user_feeds_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_feeds
    ADD CONSTRAINT user_feeds_pkey PRIMARY KEY (user_id, feed_id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: articles_feed_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX articles_feed_id_idx ON public.articles USING btree (feed_id);


--
-- Name: articles_seq_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX articles_seq_idx ON public.articles USING btree (seq);


--
-- Name: articles_site_guid_feed_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX articles_site_guid_feed_id_idx ON public.articles USING btree (site_guid, feed_id);


--
-- Name: schema_migration_version_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX schema_migration_version_idx ON public.schema_migration USING btree (version);


--
-- Name: users_auth_origin_auth_subject_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX users_auth_origin_auth_subject_idx ON public.users USING btree (auth_origin, auth_subject);


--
-- Name: articles articles_feed_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.articles
    ADD CONSTRAINT articles_feed_id_fkey FOREIGN KEY (feed_id) REFERENCES public.feeds(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--
