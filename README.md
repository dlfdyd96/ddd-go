# Implement DDD in Golang

References) [How To Implement Domain-Driven Design (DDD) in Golang](https://programmingpercy.tech/blog/how-to-domain-driven-design-ddd-golang/)

위 글을 따라하면서 DDD에 대한 간략한 정의를 정리와 함께 작은 예시를 따라하면서 배워보는 repository 입니다.

## 배경
- MSA 중요해졌어요.
- 하지만, 무분별한 사용은 복잡도를 야기시켜요.
- 위 글에서는 단계별로 `온라인 선술집` 예제를 DDD에 적용하면서 구축을 진행합니다.

## DDD?

<aside class="notice">

📌 Domain-Driven Design은 소프트웨어가 속한 도메인을 따라 소프트웨어를 구조화하고 모델링하는 방법입니다

</aside>

## Domain / Model / Ubiquitous Language / Sub-Domains


<aside class="notice">

📌 Model은 Domain을 처리하는 데 필요한 구성 요소의 추상화입니다.

</aside>

<aside class="notice">

📌 Domain은 소프트웨어가 작동할 영역입니다.

</aside>

<aside class="notice">

📌 모든 작업 참여자들(디자이너, 개발자 등등)이 사용할 공통된 언어가 Ubiquitous Language입니다.

</aside>

<aside class="notice">

📌 Sub-Domain은 루트 도메인 내부의 영역을 해결하는 데 사용되는 별도의 도메인입니다.

</aside>

## Entities & Value Objects

<aside class="notice">

📌 Entity는 변경할 수 있는 상태가 있는 고유 식별자를 참조하는 구조체입니다.

</aside>

- Entity
    - unique identifier (id)
    - Mutable (수정)

<aside class="notice">


📌 VO는 종종 도메인 내부에서 발견되며 해당 도메인의 특정 측면을 설명하는 데 사용됩니다.

</aside>

- Value Object (VO)
    - No Identifier (값 그 자체, 거래 같은 것.)
    - Immutable

## Aggregate - 결합된 Entity & VO

<aside class="notice">

📌 Aggregate는 엔터티와 값 개체가 결합된 집합입니다.

</aside>

> ***DDD Aggregate는 도메인 개념(주문, 진료소 방문, 재생 목록) — [Martin Fowler](https://martinfowler.com/bliki/DDD_Aggregate.html)***
>
- aggregate는 기본 엔터티에 대한 직접 액세스를 허용하지 않습니다.
- 예를 들어 고객과 같이 실생활에서 데이터를 올바르게 나타내기 위해 여러 엔터티가 필요한 것도 일반적입니다.
- **DDD Aggregate의 중요한 규칙은 하나의 엔터티만 루트 엔터티** 로 작동해야 한다는 것
    - 루트 엔터티의 참조가 집계를 참조하는 데에도 사용된다는 것
- Aggregate가 데이터에 대한 직접 액세스를 허용하지 않아야 하기 때문에 모든필드를 private으로 수행

## Factories - 복잡한 로직을 Encapsulate

bussiness logic을 구현할 차례.

> **Factory Pattern**은 원하는 인스턴스를 생성하는 함수에 복잡한 논리를 캡슐화하는 데 사용되는 디자인 패턴입니다. (호출자가 구현 세부 사항에 대해 알지 못하는 상태에서)

<details>
<summary>golang Elaistchsearch client example</summary>
<div markdown="1">

- `NewClient()` 함수에 `Config` 를 insert
    - `NewClient()` 함수는, elastichsearch DB에 연결하고 es document를 insert/remove 하는 client 생성하여 반환하는 Factory method 입니다.

    ```go
    func NewClient(cfg Config) (*Client, error) {
    	var addrs []string
    
    	if len(cfg.Addresses) == 0 && cfg.CloudID == "" {
    		addrs = addrsFromEnvironment()
    	} else {
    		if len(cfg.Addresses) > 0 && cfg.CloudID != "" {
    			return nil, errors.New("cannot create client: both Addresses and CloudID are set")
    		}
    
    		if cfg.CloudID != "" {
    			cloudAddr, err := addrFromCloudID(cfg.CloudID)
    			if err != nil {
    				return nil, fmt.Errorf("cannot create client: cannot parse CloudID: %s", err)
    			}
    			addrs = append(addrs, cloudAddr)
    		}
    
    		if len(cfg.Addresses) > 0 {
    			addrs = append(addrs, cfg.Addresses...)
    		}
    	}
    
    	urls, err := addrsToURLs(addrs)
    	if err != nil {
    		return nil, fmt.Errorf("cannot create client: %s", err)
    	}
    
    	if len(urls) == 0 {
    		u, _ := url.Parse(defaultURL) // errcheck exclude
    		urls = append(urls, u)
    	}
    
    	// TODO(karmi): Refactor
    	if urls[0].User != nil {
    		cfg.Username = urls[0].User.Username()
    		pw, _ := urls[0].User.Password()
    		cfg.Password = pw
    	}
    
    	tp, err := estransport.New(estransport.Config{
    		URLs:         urls,
    		Username:     cfg.Username,
    		Password:     cfg.Password,
    		APIKey:       cfg.APIKey,
    		ServiceToken: cfg.ServiceToken,
    
    		Header: cfg.Header,
    		CACert: cfg.CACert,
    
    		RetryOnStatus:        cfg.RetryOnStatus,
    		DisableRetry:         cfg.DisableRetry,
    		EnableRetryOnTimeout: cfg.EnableRetryOnTimeout,
    		MaxRetries:           cfg.MaxRetries,
    		RetryBackoff:         cfg.RetryBackoff,
    
    		CompressRequestBody: cfg.CompressRequestBody,
    
    		EnableMetrics:     cfg.EnableMetrics,
    		EnableDebugLogger: cfg.EnableDebugLogger,
    
    		DisableMetaHeader: cfg.DisableMetaHeader,
    
    		DiscoverNodesInterval: cfg.DiscoverNodesInterval,
    
    		Transport:          cfg.Transport,
    		Logger:             cfg.Logger,
    		Selector:           cfg.Selector,
    		ConnectionPoolFunc: cfg.ConnectionPoolFunc,
    	})
    	if err != nil {
    		return nil, fmt.Errorf("error creating transport: %s", err)
    	}
    
    	client := &Client{Transport: tp}
    	client.API = esapi.New(client)
    
    	if cfg.DiscoverNodesOnStart {
    		go client.DiscoverNodes()
    	}
    
    	return client, nil
    }
    ```


</div>
</details>


<aside class="notice">

📌 DDD는 복잡한 집계, 리포지토리 및 서비스를 만들기 위해 `FACTORY`를 사용할 것을 제안

</aside>

- 팩토리
    - 유효성 검사
    - ID 생성
    - 필드 초기화

## Repository - 레포지터리 패턴

<aside class="notice">

📌 DDD는 Aggregate를 저장하고 관리하는 데 `REPOSITORY`를 사용해야 합니다.

</aside>

- Interface 뒤에 스토리지/데이터베이스 솔루션의 구현을 숨기는데 의존하는 패턴
- 장점
    - 아무것도 깨지 않고 솔루션을 교환할 수 있다는 것 (OCP)
    - 테스트에도 매우 유용(단위 테스트 등을 위해 새 리포지토리를 쉽게 구현할 수 있음)
- 리포지토리를 해당 도메인과 관련된 상태로 유지해야 합니다.
    - 리포지토리를 다른 집계에 결합을 비추합니다.
    - 느슨한 결합을 위함.

## Service - 비즈니스 로직 연결

<aside class="notice">

📌 `SERVICE`는 느슨하게 연결된 모든 repository를 특정 Domain의 요구 사항을 충족하는 비즈니스 논리에 연결합니다.

</aside>

- 특정 비즈니스 논리 흐름을 수행하는 데 필요한 모든 repository를 한 서비스에서 보유함.
- 서비스 내부에 서비스가 있을 수 도 있음.

<details>
<summary>동적 팩토리를 사용하여 Service에 가변적인 configuration 을 넣는 트릭에 대하여</summary>
<div markdown="1">
```go
    type OrderConfiguration func(os *OrderService) error

    func NewOrderService(cfgs ...OrderConfiguration) (*OrderService, error) {
    	// Create the orderservice
    	os := &OrderService{}
    	// Apply all Configurations passed in
    	for _, cfg := range cfgs {
    		// Pass the service into the configuration function
    		err := cfg(os)
    		if err != nil {
    			return nil, err
    		}
    	}
    	return os, nil
    }
    ```

이 코드의 핵심은 OrderConfiguration 라는 alias function을 매개변수로 받는 NewOrderService 입니다.

동적 팩토리를 허용하여 가변적인 OrderConfiguration 들을 받을 수 있도록 구성합니다.

이 방법은 서비스의 특정 부분을 원하는 repository로 교체할 수 있으므로 **단위 테스트에 매우 유용**합니다.

작은 서비스에 과도하게 보일 수 있지만, 더나아가 서비스의 내부 설정 및 옵션에도 사용할 수 있다는 점이 있습니다.

그래서 이 방법(동적팩토리)를 통해 OrderService에 customerRepository를 어떻게 주입하냐!

    ```go
    // WithCustomerRepository applies a given customer repository to the OrderService
    func WithCustomerRepository(cr customer.CustomerRepository) OrderConfiguration {
    	// return a function that matches the OrderConfiguration alias,
    	// You need to return this so that the parent function can take in all the needed parameters
    	return func(os *OrderService) error {
    		os.customers = cr
    		return nil
    	}
    }
    
    // WithMemoryCustomerRepository applies a memory customer repository to the OrderService
    func WithMemoryCustomerRepository() OrderConfiguration {
    	// Create the memory repo, if we needed parameters, such as connection strings they could be inputted here
    	cr := memory.New()
    	return WithCustomerRepository(cr)
    }
    ```

위의 함수를 사용하여

    ```go
    // In Memory Example used in Development
    NewOrderService(WithMemoryCustomerRepository())
    // We could in the future switch to MongoDB like this
    NewOrderService(WithMongoCustomerRepository())
    
    // =========== or ===
    
    products := []Product{prod1, prod2, prod3}
    os, err := NewOrderService(
    		WithMemoryCustomerRepository(),
    		WithMemoryProductRepository(products),
    	)
    ```

다음과 같이 쓸 수 있습니다.
</div>
</details>

## 정리

- **Entity** - 식별 가능한 가변 구조체.
- **Value Object** — 불변의 식별 불가능한 구조체.
- **Aggregates** — 리포지토리에 저장된 엔티티 및 값 개체의 결합된 집합입니다.
- **Repository** — 집계 또는 기타 정보를 저장하는 구현
- **Factory** — 복잡한 개체를 만들고 다른 도메인의 개발자가 새 인스턴스를 쉽게 만들 수 있도록 하는 생성자
- **Service** — 비즈니스 흐름을 함께 구축하는 리포지토리 및 하위 서비스 모음