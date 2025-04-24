package utils

import (
	"encoding/json"
	"errors"
	"log"
	"slices"
)

type Movie struct {
	Adult           bool    `json:"adult"`
	BackdropPath    string  `json:"backdrop_path"`
	GenreIds        []int   `json:"genre_ids"`
	Id              int     `json:"id"`
	OriginalLanguage string  `json:"original_language"`
	OriginalTitle   string  `json:"original_title"`
	Overview        string  `json:"overview"`
	Popularity      float64 `json:"popularity"`
	Poster          string  `json:"poster_path"`
	ReleaseDate     string  `json:"release_date"`
	Title           string  `json:"title"`
	Video           bool    `json:"video"`
	VoteAverage     float64 `json:"vote_average"`
	VoteCount       int     `json:"vote_count"`
}

type MovieResponse struct {
	Page         int `json:"page"`
	Results      []Movie `json:"results"`
	TotalPages   int `json:"total_pages"`
	TotalResults int `json:"total_results"`
}

type RecommendationResponse struct {
	Page         int `json:"page"`
	Results      []Movie `json:"results"`
	TotalPages   int `json:"total_pages"`
	TotalResults int `json:"total_results"`
}

func MovieMeetsUsageCriteria(movie Movie) bool {
	return movie.Poster != "" && movie.Popularity > 1 && movie.VoteAverage > 2 && movie.VoteCount > 25
}

func SearchMovies(search string) (MovieResponse, error) {
	str := `{
		"page": 1,
		"results": [
			{
			"adult": false,
			"backdrop_path": "/hdHIjZxq3SWFqpAz4NFhdbud0iz.jpg",
			"genre_ids": [
				27,
				878
			],
			"id": 348,
			"original_language": "en",
			"original_title": "Alien",
			"overview": "During its return to the earth, commercial spaceship Nostromo intercepts a distress signal from a distant planet. When a three-member team of the crew discovers a chamber containing thousands of eggs on the planet, a creature inside one of the eggs attacks an explorer. The entire crew is unaware of the impending nightmare set to descend upon them when the alien parasite planted inside its unfortunate host is birthed.",
			"popularity": 48.9902,
			"poster_path": "/vfrQk5IPloGg1v9Rzbh2Eg3VGyM.jpg",
			"release_date": "1979-05-25",
			"title": "Alien",
			"video": false,
			"vote_average": 8.2,
			"vote_count": 15216
			}
		],
		"total_pages": 58,
		"total_results": 1154
	}`

	var response MovieResponse
	err := json.Unmarshal([]byte(str), &response)
	if err != nil {
		return MovieResponse{}, err
	}

	// remove movies with no poster, low popularity, low vote average, low vote count, no video, or is adult
	// TODO: might look into filtering these out at the request level
	filteredResults := []Movie{}
	for _, movie := range response.Results {
		log.Printf("movie: %+v", movie)
		if MovieMeetsUsageCriteria(movie) { 
			filteredResults = append(filteredResults, movie)
		}
	}

	response.Results = filteredResults
	return response, nil
} 


func GetMovieRecommendations(id string) (RecommendationResponse, error) {
	recommendationStr := `{
		"page": 1,
		"results": [
			{
			"backdrop_path": "/jMBpJFRtrtIXymer93XLavPwI3P.jpg",
			"id": 679,
			"title": "Aliens",
			"original_title": "Aliens",
			"overview": "Ripley, the sole survivor of the Nostromo's deadly encounter with the monstrous Alien, returns to Earth after drifting through space in hypersleep for 57 years. Although her story is initially met with skepticism, she agrees to accompany a team of Colonial Marines back to LV-426.",
			"poster_path": "/r1x5JGpyqZU8PYhbs4UcrO1Xb6x.jpg",
			"media_type": "movie",
			"adult": false,
			"original_language": "en",
			"genre_ids": [
				28,
				53,
				878
			],
			"popularity": 16.5073,
			"release_date": "1986-07-18",
			"video": false,
			"vote_average": 7.951,
			"vote_count": 10059
			},
			{
			"backdrop_path": "/nEmOmbCWBXS3tHU2N49z693KDK.jpg",
			"id": 8077,
			"title": "Alien³",
			"original_title": "Alien³",
			"overview": "After escaping with Newt and Hicks from the alien planet, Ripley crash lands on Fiorina 161, a prison planet and host to a correctional facility. Unfortunately, although Newt and Hicks do not survive the crash, a more unwelcome visitor does. The prison does not allow weapons of any kind, and with aid being a long time away, the prisoners must simply survive in any way they can.",
			"poster_path": "/xh5wI0UoW7DfS1IyLy3d2CgrCEP.jpg",
			"media_type": "movie",
			"adult": false,
			"original_language": "en",
			"genre_ids": [
				878,
				28,
				27
			],
			"popularity": 8.3829,
			"release_date": "1992-05-22",
			"video": false,
			"vote_average": 6.37,
			"vote_count": 5869
			},
			{
			"backdrop_path": "/ikr0UILfvRerzMNoBTtJtyuWAEV.jpg",
			"id": 8078,
			"title": "Alien Resurrection",
			"original_title": "Alien Resurrection",
			"overview": "Two hundred years after Lt. Ripley died, a group of scientists clone her, hoping to breed the ultimate weapon. But the new Ripley is full of surprises … as are the new aliens. Ripley must team with a band of smugglers to keep the creatures from reaching Earth.",
			"poster_path": "/9aRDMlU5Zwpysilm0WCWzU2PCFv.jpg",
			"media_type": "movie",
			"adult": false,
			"original_language": "en",
			"genre_ids": [
				878,
				27,
				28
			],
			"popularity": 10.9485,
			"release_date": "1997-11-12",
			"video": false,
			"vote_average": 6.159,
			"vote_count": 5112
			},
			{
			"backdrop_path": "/3nYlM34QhzdtAvWRV5bN4nLtnTc.jpg",
			"id": 578,
			"title": "Jaws",
			"original_title": "Jaws",
			"overview": "When the seaside community of Amity finds itself under attack by a dangerous great white shark, the town's chief of police, a young marine biologist, and a grizzled hunter embark on a desperate quest to destroy the beast before it strikes again.",
			"poster_path": "/lxM6kqilAdpdhqUl2biYp5frUxE.jpg",
			"media_type": "movie",
			"adult": false,
			"original_language": "en",
			"genre_ids": [
				27,
				53,
				12
			],
			"popularity": 15.9236,
			"release_date": "1975-06-20",
			"video": false,
			"vote_average": 7.664,
			"vote_count": 10714
			},
			{
			"backdrop_path": "/qr7dUqleMRd0VgollazbmyP9XjI.jpg",
			"id": 78,
			"title": "Blade Runner",
			"original_title": "Blade Runner",
			"overview": "In the smog-choked dystopian Los Angeles of 2019, blade runner Rick Deckard is called out of retirement to terminate a quartet of replicants who have escaped to Earth seeking their creator for a way to extend their short life spans.",
			"poster_path": "/63N9uy8nd9j7Eog2axPQ8lbr3Wj.jpg",
			"media_type": "movie",
			"adult": false,
			"original_language": "en",
			"genre_ids": [
				878,
				18,
				53
			],
			"popularity": 16.1026,
			"release_date": "1982-06-25",
			"video": false,
			"vote_average": 7.941,
			"vote_count": 14080
			},
			{
			"backdrop_path": "/2qluV8y79LnBBHaMpwewCjQ1Htk.jpg",
			"id": 126889,
			"title": "Alien: Covenant",
			"original_title": "Alien: Covenant",
			"overview": "The crew of the colony ship Covenant, bound for a remote planet on the far side of the galaxy, discovers what they think is an uncharted paradise but is actually a dark, dangerous world.",
			"poster_path": "/zecMELPbU5YMQpC81Z8ImaaXuf9.jpg",
			"media_type": "movie",
			"adult": false,
			"original_language": "en",
			"genre_ids": [
				878,
				27,
				53
			],
			"popularity": 10.5279,
			"release_date": "2017-05-09",
			"video": false,
			"vote_average": 6.2,
			"vote_count": 8855
			},
			{
			"backdrop_path": "/9OKzKuCfcMwpHVp56pjNI7J4xXR.jpg",
			"id": 923,
			"title": "Dawn of the Dead",
			"original_title": "Dawn of the Dead",
			"overview": "During an ever-growing epidemic of zombies that have risen from the dead, two Philadelphia SWAT team members, a traffic reporter, and his television-executive girlfriend seek refuge in a secluded shopping mall.",
			"poster_path": "/3gx7VwU0OmvMihtiDmRNmsd4OG5.jpg",
			"media_type": "movie",
			"adult": false,
			"original_language": "en",
			"genre_ids": [
				27
			],
			"popularity": 7.1993,
			"release_date": "1978-09-02",
			"video": false,
			"vote_average": 7.5,
			"vote_count": 2131
			},
			{
			"backdrop_path": "/mhgfOLtAiJk1nm3kKlkiUKl7e9M.jpg",
			"id": 2252,
			"title": "Eastern Promises",
			"original_title": "Eastern Promises",
			"overview": "A Russian teenager living in London dies during childbirth but leaves clues in her diary that could tie her child to a rape involving a violent Russian mob family.",
			"poster_path": "/dpiJWb4NrWgcOg2rusuLhDM0hTm.jpg",
			"media_type": "movie",
			"adult": false,
			"original_language": "en",
			"genre_ids": [
				53,
				80,
				9648
			],
			"popularity": 3.8609,
			"release_date": "2007-09-14",
			"video": false,
			"vote_average": 7.355,
			"vote_count": 3444
			},
			{
			"backdrop_path": "/r9leYNa8nTRCceZrZhP1DXkgKVb.jpg",
			"id": 1091,
			"title": "The Thing",
			"original_title": "The Thing",
			"overview": "In the winter of 1982, a twelve-man research team at a remote Antarctic research station discovers an alien buried in the snow for over 100,000 years. Soon unfrozen, the form-changing creature wreaks havoc, creates terror... and becomes one of them.",
			"poster_path": "/tzGY49kseSE9QAKk47uuDGwnSCu.jpg",
			"media_type": "movie",
			"adult": false,
			"original_language": "en",
			"genre_ids": [
				27,
				9648,
				878
			],
			"popularity": 10.4936,
			"release_date": "1982-06-25",
			"video": false,
			"vote_average": 8.068,
			"vote_count": 7249
			},
			{
			"backdrop_path": "/6ghc2ySuEnTW8h3YgZZufk2SOTv.jpg",
			"id": 12113,
			"title": "Body of Lies",
			"original_title": "Body of Lies",
			"overview": "The CIA's hunt is on for the mastermind of a wave of terrorist attacks. Roger Ferris is the agency's man on the ground, moving from place to place, scrambling to stay ahead of ever-shifting events. An eye in the sky – a satellite link – watches Ferris.  At the other end of that real-time link is the CIA's Ed Hoffman, strategizing events from thousands of miles away. And as Ferris nears the target, he discovers trust can be just as dangerous as it is necessary for survival.",
			"poster_path": "/rNEZug6er0bIj9LVN2JaQig6oZy.jpg",
			"media_type": "movie",
			"adult": false,
			"original_language": "en",
			"genre_ids": [
				28,
				18,
				53
			],
			"popularity": 3.2426,
			"release_date": "2008-10-09",
			"video": false,
			"vote_average": 6.629,
			"vote_count": 3092
			},
			{
			"backdrop_path": "/w5IDXtifKntw0ajv2co7jFlTQDM.jpg",
			"id": 62,
			"title": "2001: A Space Odyssey",
			"original_title": "2001: A Space Odyssey",
			"overview": "Humanity finds a mysterious object buried beneath the lunar surface and sets off to find its origins with the help of HAL 9000, the world's most advanced super computer.",
			"poster_path": "/ve72VxNqjGM69Uky4WTo2bK6rfq.jpg",
			"media_type": "movie",
			"adult": false,
			"original_language": "en",
			"genre_ids": [
				878,
				9648,
				12
			],
			"popularity": 13.2731,
			"release_date": "1968-04-02",
			"video": false,
			"vote_average": 8.1,
			"vote_count": 11795
			},
			{
			"backdrop_path": "/m9yHpYz4GcEtmHJW4rvQIrF891h.jpg",
			"id": 534,
			"title": "Terminator Salvation",
			"original_title": "Terminator Salvation",
			"overview": "All grown up in post-apocalyptic 2018, John Connor must lead the resistance of humans against the increasingly dominating militaristic robots. But when Marcus Wright appears, his existence confuses the mission as Connor tries to determine whether Wright has come from the future or the past -- and whether he's friend or foe.",
			"poster_path": "/gw6JhlekZgtKUFlDTezq3j5JEPK.jpg",
			"media_type": "movie",
			"adult": false,
			"original_language": "en",
			"genre_ids": [
				28,
				878,
				53
			],
			"popularity": 11.7307,
			"release_date": "2009-05-20",
			"video": false,
			"vote_average": 6.06,
			"vote_count": 6588
			},
			{
			"backdrop_path": "/qUq3QTr2KLvGIcN0GaaaYx9bbyH.jpg",
			"id": 510,
			"title": "One Flew Over the Cuckoo's Nest",
			"original_title": "One Flew Over the Cuckoo's Nest",
			"overview": "A petty criminal fakes insanity to serve his sentence in a mental ward rather than prison. He soon finds himself as a leader to the other patients—and an enemy to the cruel, domineering nurse who runs the ward.",
			"poster_path": "/kjWsMh72V6d8KRLV4EOoSJLT1H7.jpg",
			"media_type": "movie",
			"adult": false,
			"original_language": "en",
			"genre_ids": [
				18
			],
			"popularity": 18.5782,
			"release_date": "1975-11-19",
			"video": false,
			"vote_average": 8.414,
			"vote_count": 10728
			},
			{
			"backdrop_path": "/9Qs9oyn4iE8QtQjGZ0Hp2WyYNXT.jpg",
			"id": 28,
			"title": "Apocalypse Now",
			"original_title": "Apocalypse Now",
			"overview": "At the height of the Vietnam war, Captain Benjamin Willard is sent on a dangerous mission that, officially, \"does not exist, nor will it ever exist.\" His goal is to locate - and eliminate - a mysterious Green Beret Colonel named Walter Kurtz, who has been leading his personal army on illegal guerrilla missions into enemy territory.",
			"poster_path": "/gQB8Y5RCMkv2zwzFHbUJX3kAhvA.jpg",
			"media_type": "movie",
			"adult": false,
			"original_language": "en",
			"genre_ids": [
				18,
				10752
			],
			"popularity": 10.0738,
			"release_date": "1979-05-19",
			"video": false,
			"vote_average": 8.272,
			"vote_count": 8462
			},
			{
			"backdrop_path": "/1KgXxv6tHXOnakqYvMPvFwYKWiw.jpg",
			"id": 762,
			"title": "Monty Python and the Holy Grail",
			"original_title": "Monty Python and the Holy Grail",
			"overview": "King Arthur, accompanied by his squire, recruits his Knights of the Round Table, including Sir Bedevere the Wise, Sir Lancelot the Brave, Sir Robin the Not-Quite-So-Brave-As-Sir-Lancelot and Sir Galahad the Pure. On the way, Arthur battles the Black Knight who, despite having had all his limbs chopped off, insists he can still fight. They reach Camelot, but Arthur decides not  to enter, as \"it is a silly place\".",
			"poster_path": "/8AVb7tyxZRsbKJNOTJHQZl7JYWO.jpg",
			"media_type": "movie",
			"adult": false,
			"original_language": "en",
			"genre_ids": [
				12,
				35,
				14
			],
			"popularity": 9.6163,
			"release_date": "1975-03-14",
			"video": false,
			"vote_average": 7.802,
			"vote_count": 5956
			},
			{
			"backdrop_path": "/ahUaAgnkFu7QlBh5h4LCNeaSurV.jpg",
			"id": 218,
			"title": "The Terminator",
			"original_title": "The Terminator",
			"overview": "In the post-apocalyptic future, reigning tyrannical supercomputers teleport a cyborg assassin known as the \"Terminator\" back to 1984 to kill Sarah Connor, whose unborn son is destined to lead insurgents against 21st century mechanical hegemony. Meanwhile, the human-resistance movement dispatches a lone warrior to safeguard Sarah. Can he stop the virtually indestructible killing machine?",
			"poster_path": "/hzXSE66v6KthZ8nPoLZmsi2G05j.jpg",
			"media_type": "movie",
			"adult": false,
			"original_language": "en",
			"genre_ids": [
				28,
				53,
				878
			],
			"popularity": 44.0515,
			"release_date": "1984-10-26",
			"video": false,
			"vote_average": 7.664,
			"vote_count": 13524
			},
			{
			"backdrop_path": "/olz8Xw3yOLpBAHKgPoSRwmomdM.jpg",
			"id": 13448,
			"title": "Angels & Demons",
			"original_title": "Angels & Demons",
			"overview": "Harvard symbologist Robert Langdon is recruited by the Vatican to investigate the apparent return of the Illuminati – a secret, underground organization – after four cardinals are kidnapped on the night of the papal conclave.",
			"poster_path": "/tFZQAuulEOtFTp0gHbVdEXwGrYe.jpg",
			"media_type": "movie",
			"adult": false,
			"original_language": "en",
			"genre_ids": [
				53,
				9648
			],
			"popularity": 20.2611,
			"release_date": "2009-04-23",
			"video": false,
			"vote_average": 6.715,
			"vote_count": 6911
			},
			{
			"backdrop_path": "/aRka9neADW1M0Zf9lF8kW2jEgXe.jpg",
			"id": 948,
			"title": "Halloween",
			"original_title": "Halloween",
			"overview": "Fifteen years after murdering his sister on Halloween Night 1963, Michael Myers escapes from a mental hospital and returns to the small town of Haddonfield, Illinois to kill again.",
			"poster_path": "/wijlZ3HaYMvlDTPqJoTCWKFkCPU.jpg",
			"media_type": "movie",
			"adult": false,
			"original_language": "en",
			"genre_ids": [
				27,
				53
			],
			"popularity": 13.7624,
			"release_date": "1978-10-24",
			"video": false,
			"vote_average": 7.562,
			"vote_count": 5754
			},
			{
			"backdrop_path": "/vINgGecnz95iDL6fjQMARDsocgG.jpg",
			"id": 280,
			"title": "Terminator 2: Judgment Day",
			"original_title": "Terminator 2: Judgment Day",
			"overview": "Set ten years after the events of the original, James Cameron's classic sci-fi action flick tells the story of a second attempt to get rid of rebellion leader John Connor, this time targeting the boy himself. However, the rebellion has sent a reprogrammed terminator to protect Connor.",
			"poster_path": "/5M0j0B18abtBI5gi2RhfjjurTqb.jpg",
			"media_type": "movie",
			"adult": false,
			"original_language": "en",
			"genre_ids": [
				28,
				53,
				878
			],
			"popularity": 16.3274,
			"release_date": "1991-07-03",
			"video": false,
			"vote_average": 8.124,
			"vote_count": 13251
			},
			{
			"backdrop_path": "/jVHRmhSMlcEhTpYzSBgd6QvFHQA.jpg",
			"id": 395,
			"title": "AVP: Alien vs. Predator",
			"original_title": "AVP: Alien vs. Predator",
			"overview": "When scientists discover something near Antarctica that appears to be a buried Pyramid, they send a research team out to investigate. Little do they know that they are about to step into a hunting ground where Aliens are grown as sport for the Predator race.",
			"poster_path": "/ySWu5bCnnmgV1cVacvFnFIhgOjp.jpg",
			"media_type": "movie",
			"adult": false,
			"original_language": "en",
			"genre_ids": [
				12,
				878,
				28,
				27
			],
			"popularity": 7.1216,
			"release_date": "2004-08-12",
			"video": false,
			"vote_average": 5.9,
			"vote_count": 4551
			}
		],
		"total_pages": 2,
		"total_results": 40
	}`

	recommendations := RecommendationResponse{}
	json.Unmarshal([]byte(recommendationStr), &recommendations)

	filteredResults := []Movie{}
	for _, movie := range recommendations.Results {
		if MovieMeetsUsageCriteria(movie) { 
			filteredResults = append(filteredResults, movie)
		}
	}

	recommendations.Results = filteredResults

	return recommendations, nil
}

type Trailer struct {
	ISO6391 string `json:"iso_639_1"`
	ISO31661 string `json:"iso_3166_1"`
	Name string `json:"name"`
	Key string `json:"key"`
	Site string `json:"site"`
	Size int `json:"size"`
	Type string `json:"type"`
	Official bool `json:"official"`
	PublishedAt string `json:"published_at"`
	Id string `json:"id"`
}
	
func GetBestMovieTrailer(id string) (Trailer, error) {
	trailerStr := `[
		{
		"iso_639_1": "en",
		"iso_3166_1": "US",
		"name": "Modern Trailer",
		"key": "sVwH0hIvV5k",
		"site": "YouTube",
		"size": 1080,
		"type": "Trailer",
		"official": true,
		"published_at": "2020-10-16T19:00:06.000Z",
		"id": "5f8e964ece9e9100334cac29"
		},
		{
		"iso_639_1": "en",
		"iso_3166_1": "US",
		"name": "Tom Skerritt on the Ripley character in ALIEN",
		"key": "MdIQBsIN-iw",
		"site": "YouTube",
		"size": 1080,
		"type": "Featurette",
		"official": true,
		"published_at": "2020-07-23T19:51:14.000Z",
		"id": "65a42a635a788400c6c968a0"
		}
	]`

	trailers := []Trailer{}
	json.Unmarshal([]byte(trailerStr), &trailers)

	filteredTrailers := []Trailer{}
	for _, trailer := range trailers {
		if trailer.Site == "YouTube" {
			filteredTrailers = append(filteredTrailers, trailer)
		}
	}

	slices.SortFunc(filteredTrailers, func(a, b Trailer) int {
		if a.Type == "Trailer" && b.Type != "Trailer" {
			return -1;
		} else if a.Type != "Trailer" && b.Type == "Trailer" {
			return 1;
		} else if a.Official && !b.Official {
			return -1;
		} else if !a.Official && b.Official {
			return 1;
		} else {
			return 0;
		}
	})

	if len(filteredTrailers) > 0 {
		return filteredTrailers[0], nil
	}

	return Trailer{}, errors.New("no trailers found")
}

